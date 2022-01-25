package paxos

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/fudute/paxos/config"
	pb "github.com/fudute/paxos/protoc"
)

type IDService interface {
	Next() int64
}

type Proposer struct {
	self              string
	quorum            int
	proposalGenerator IDService
	cluster           *config.Cluster
	conns             *ConnectionsManager

	mu               sync.Mutex
	log              []slot
	smallestUnchosen int64

	// leader 信息，在leader信息有效时间内，只会接受来自leader的prepare和accept请求
	leader   string
	expire   *time.Timer
	interval time.Duration

	pb.UnimplementedProposerServer
	pb.UnimplementedLearnerServer
}

type status int

const (
	idle    status = 0
	running status = 1
	chosen  status = 2
)

type slot struct {
	value  string
	status status
}

func NewProposer(self string, quorum int, cluster *config.Cluster, generator IDService) *Proposer {
	proposer := &Proposer{
		self:              self,
		quorum:            quorum,
		proposalGenerator: generator,
		interval:          time.Second,
		cluster:           cluster,
		conns:             NewConnectionsManager(),
	}
	return proposer
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
		cli, err := p.conns.Proposer(p.leader)
		if err == nil {
			return cli.Propose(ctx, req)
		}
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
		p.log = append(p.log, make([]slot, len(p.log)+1)...)
	}

	p.log[index].value = req.ProposalValue
	p.log[index].status = chosen

	for p.smallestUnchosen < int64(len(p.log)) && p.log[p.smallestUnchosen].status == chosen {
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

func (p *Proposer) pickSlot() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	index := p.smallestUnchosen

	for index < int64(len(p.log)) && p.log[index].status != idle {
		index++
	}

	if index == int64(len(p.log)) {
		p.log = append(p.log, slot{})
	}
	p.log[index].status = running
	return index
}

func (p *Proposer) chooseOne(value string) (chosen string) {
	index := p.pickSlot()
	for {
		proposalNum := p.proposalGenerator.Next()
		prepareReplys := p.broadcast(p.cluster.Acceptors, func(node *config.Node) interface{} { return p.onPrepare(index, proposalNum, node) })

		var mini int64
		for i := 0; i <= p.quorum; i++ {
			reply := (<-prepareReplys).(*pb.PrepareReply)
			if reply.AcceptedProposal > mini && reply.AcceptedValue != "" {
				value = reply.AcceptedValue
				mini = reply.AcceptedProposal
			}
		}

		acceptReplys := p.broadcast(p.cluster.Acceptors, func(node *config.Node) interface{} { return p.onAccept(index, proposalNum, value, node) })

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
			p.broadcast(p.cluster.Learners, func(node *config.Node) interface{} {
				return p.OnLearn(index, value, node)
			})
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

func (p *Proposer) onPrepare(index, proposalNum int64, node *config.Node) *pb.PrepareReply {
	cli, err := p.conns.Acceptor(node.Addr)
	if err != nil {
		return nil
	}
	reply, err := cli.Prepare(context.Background(), &pb.PrepareRequest{
		Index:       index,
		ProposalNum: proposalNum,
	})
	if err != nil {
		log.Printf("[P%s] Prepare on acceptor %v failed %v", p.self, node.Addr, err)
		return nil
	}
	return reply
}

func (p *Proposer) onAccept(index, proposalNum int64, value string, node *config.Node) *pb.AcceptReply {
	cli, err := p.conns.Acceptor(node.Addr)
	if err != nil {
		return nil
	}
	reply, err := cli.Accept(context.Background(), &pb.AcceptRequest{
		Index:         index,
		ProposalNum:   proposalNum,
		ProposalValue: value,
	})
	if err != nil {
		log.Printf("[P%s] Accept on acceptor %v failed %v", p.self, node.Addr, err)
		return nil
	}
	return reply
}

func (p *Proposer) OnLearn(index int64, value string, node *config.Node) *pb.LearnReply {
	cli, err := p.conns.Learner(node.Addr)
	if err != nil {
		return nil
	}
	reply, err := cli.Learn(context.Background(), &pb.LearnRequest{
		Index:         index,
		Proposer:      p.self,
		ProposalValue: value,
	})
	if err != nil {
		log.Printf("[P%s] Learn on learner %v failed %v", p.self, node.Addr, err)
		return nil
	}
	return reply
}
