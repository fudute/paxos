package paxos

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fudute/paxos/config"
	pb "github.com/fudute/paxos/protoc"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type IDService interface {
	Next() int64
}

type Proposer struct {
	self              string
	quorum            int
	proposalGenerator IDService
	acceptors         []*config.Node
	acceptorsClient   map[string]pb.AcceptorClient
	proposers         []*config.Node
	proposersClient   map[string]pb.ProposerClient

	mu               sync.Mutex
	log              []string
	smallestUnchosen int64

	// leader 信息，在leader信息有效时间内，只会接受来自leader的prepare和accept请求
	leader   string
	expire   *time.Timer
	interval time.Duration

	pb.UnimplementedProposerServer
	pb.UnimplementedLeanerServer
}

func NewProposer(self string, quorum int, cluster *config.Cluster, generator IDService) *Proposer {
	proposer := &Proposer{
		self:              self,
		quorum:            quorum,
		acceptors:         cluster.Acceptors,
		proposers:         cluster.Proposers,
		acceptorsClient:   make(map[string]pb.AcceptorClient),
		proposersClient:   make(map[string]pb.ProposerClient),
		proposalGenerator: generator,
		interval:          time.Second,
	}
	for _, acceptor := range cluster.Acceptors {
		cc, _ := grpc.Dial(acceptor.Addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.DefaultConfig,
			}),
		)
		proposer.acceptorsClient[acceptor.Addr] = pb.NewAcceptorClient(cc)
	}

	for _, p := range cluster.Proposers {
		cc, _ := grpc.Dial(p.Addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.DefaultConfig,
			}),
		)
		proposer.proposersClient[p.Addr] = pb.NewProposerClient(cc)
	}

	return proposer
}

func (p *Proposer) String() string {
	b := strings.Builder{}
	for _, s := range p.log {
		b.WriteString(s)
		b.WriteByte(' ')
	}
	return b.String()
}

func (p *Proposer) LeaderExpired() bool {
	if p.leader == "" {
		return true
	}
	select {
	case <-p.expire.C:
		return true
	default:
		return false
	}
}

func (p *Proposer) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeReply, error) {
	if p.leader != p.self && !p.LeaderExpired() {
		return p.proposersClient[p.leader].Propose(ctx, req)
	}
	for {
		chosen := p.chooseOne(req.ProposalValue)
		if req.ProposalValue == chosen {
			break
		}
	}
	return &pb.ProposeReply{}, nil
}

func (p *Proposer) Learn(ctx context.Context, req *pb.LearnRequest) (*pb.LearnReply, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	index := req.Index
	if index >= int64(len(p.log)) {
		p.log = append(p.log, make([]string, len(p.log)+1)...)
	}

	p.log[index] = req.ProposalValue

	for p.smallestUnchosen < int64(len(p.log)) && p.log[p.smallestUnchosen] != "" {
		p.smallestUnchosen++
	}

	if p.leader != req.Proposer {
		log.Printf("[P%v] leader change from %v to %v", p.self, p.leader, req.Proposer)
	}
	p.leader = req.Proposer
	if p.expire == nil {
		p.expire = time.NewTimer(p.interval)
	} else {
		p.expire.Reset(p.interval)
	}

	return &pb.LearnReply{}, nil
}

func (p *Proposer) returnFirstUnChosen() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	for p.smallestUnchosen < int64(len(p.log)) && p.log[p.smallestUnchosen] != "" {
		p.smallestUnchosen++
	}
	return p.smallestUnchosen
}

func (p *Proposer) chooseOne(value string) (chosen string) {
	index := p.returnFirstUnChosen()
	for {
		proposalNum := p.proposalGenerator.Next()
		prepareReplys := p.broadcast(p.acceptors, func(node *config.Node) interface{} { return p.onPrepare(index, proposalNum, node.Addr) })

		var mini int64
		for i := 0; i <= p.quorum; i++ {
			reply := (<-prepareReplys).(*pb.PrepareReply)
			if reply.AcceptedProposal > mini && reply.AcceptedValue != "" {
				value = reply.AcceptedValue
				mini = reply.AcceptedProposal
			}
		}

		acceptReplys := p.broadcast(p.acceptors, func(node *config.Node) interface{} { return p.onAccept(index, proposalNum, node.Addr, value) })

		isChosen := true
		for i := 0; i <= p.quorum; i++ {
			reply := (<-acceptReplys).(*pb.AcceptReply)
			if reply.GetMiniProposal() > proposalNum {
				isChosen = false
				break
			}
		}

		if isChosen {
			log.Printf("[P%s] value at index %v is chosen: %v ", p.self, index, value)
			p.log[index] = value
			// go func() {
			// 	for _, proposer := range p.proposers {
			// 		proposer.Learn(context.Background(), &pb.LearnRequest{
			// 			Index:         index,
			// 			ProposalValue: value,
			// 			Proposer:      p.self,
			// 		})
			// 	}
			// }()
			break
		}
	}
	return value
}

func (p *Proposer) broadcast(nodes []*config.Node, f func(node *config.Node) interface{}) <-chan interface{} {
	out := make(chan interface{})

	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(node *config.Node) {
			defer wg.Done()
			reply := f(node)
			if reply != nil {
				out <- reply
			}
		}(node)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (p *Proposer) onPrepare(index, proposalNum int64, addr string) *pb.PrepareReply {
	if cli, ok := p.acceptorsClient[addr]; ok {
		reply, err := cli.Prepare(context.Background(), &pb.PrepareRequest{
			Index:       index,
			ProposalNum: proposalNum,
		})
		if err != nil {
			log.Printf("[P%s] Prepare on acceptor %v failed %v", p.self, addr, err)
			return nil
		}
		return reply
	}
	return nil
}

func (p *Proposer) onAccept(index, proposalNum int64, addr string, value string) *pb.AcceptReply {
	if cli, ok := p.acceptorsClient[addr]; ok {
		reply, err := cli.Accept(context.Background(), &pb.AcceptRequest{
			Index:         index,
			ProposalNum:   proposalNum,
			ProposalValue: value,
		})
		if err != nil {
			log.Printf("[P%s] Accept on acceptor %v failed %v", p.self, addr, err)
			return nil
		}
		return reply
	}
	return nil
}
