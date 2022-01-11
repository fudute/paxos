package paxos

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

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
	acceptors         map[string]pb.AcceptorClient
	proposers         map[string]pb.ProposerClient

	mu               sync.Mutex
	log              []string
	smallestUnchosen int64

	// leader 信息，在leader信息有效时间内，只会接受来自leader的prepare和accept请求
	leader   string
	expire   *time.Timer
	interval time.Duration

	pb.UnimplementedProposerServer
}

func NewProposer(self string, quorum int, acceptors []string, proposers []string, generator IDService) *Proposer {
	proposer := &Proposer{
		self:              self,
		quorum:            quorum,
		acceptors:         make(map[string]pb.AcceptorClient),
		proposers:         make(map[string]pb.ProposerClient),
		proposalGenerator: generator,
		interval:          time.Second,
	}
	for _, acceptor := range acceptors {
		cc, _ := grpc.Dial(acceptor,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.DefaultConfig,
			}),
		)
		proposer.acceptors[acceptor] = pb.NewAcceptorClient(cc)
	}

	for _, p := range proposers {
		cc, _ := grpc.Dial(p,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.DefaultConfig,
			}),
		)
		proposer.proposers[p] = pb.NewProposerClient(cc)
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
		return false
	}
	select {
	case <-p.expire.C:
		return true
	default:
		return false
	}
}

func (p *Proposer) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeReply, error) {
	if p.leader != "" && p.leader != p.self && !p.LeaderExpired() {
		reply, err := p.proposers[p.leader].Propose(ctx, req)
		if err != nil {
			log.Printf("[P%v] redirect request to leader %v failed", p.self, p.leader)
		}
		return reply, err
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
		preAlloc := make([]string, index-int64(len(p.log))+1)
		p.log = append(p.log, preAlloc...)
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
		prepareReplys := p.broadcast(proposalNum, func(n int64, addr string) interface{} { return p.onPrepare(index, n, addr) })

		var mini int64
		for i := 0; i <= p.quorum; i++ {
			reply := (<-prepareReplys).(*pb.PrepareReply)
			if reply.AcceptedProposal > mini && reply.AcceptedValue != "" {
				value = reply.AcceptedValue
				mini = reply.AcceptedProposal
			}
		}

		acceptReplys := p.broadcast(proposalNum, func(n int64, addr string) interface{} { return p.onAccept(index, n, addr, value) })

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
			go func() {
				for _, proposer := range p.proposers {
					proposer.Learn(context.Background(), &pb.LearnRequest{
						Index:         index,
						ProposalValue: value,
						Proposer:      p.self,
					})
				}
			}()
			break
		}
	}
	return value
}

func (p *Proposer) broadcast(n int64, f func(n int64, addr string) interface{}) <-chan interface{} {
	out := make(chan interface{})

	var wg sync.WaitGroup
	for addr, acceptor := range p.acceptors {
		wg.Add(1)
		go func(addr string, acceptor pb.AcceptorClient) {
			defer wg.Done()
			reply := f(n, addr)
			if reply != nil {
				out <- reply
			}
		}(addr, acceptor)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (p *Proposer) onPrepare(index, proposalNum int64, addr string) *pb.PrepareReply {
	if cli, ok := p.acceptors[addr]; ok {
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
	if cli, ok := p.acceptors[addr]; ok {
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
