package paxos

import (
	"context"
	"log"
	"strings"
	"sync"

	pb "github.com/fudute/paxos/protoc"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type IDService interface {
	Next() int64
}

type Proposer struct {
	id                int
	quorum            int
	proposalGenerator IDService
	acceptors         map[string]pb.PaxosClient

	mu                 sync.Mutex
	log                []*Slot
	largestChosenIndex int64

	wg sync.WaitGroup
}

type SlotStatus int

const (
	Unused    SlotStatus = 0
	Chooseing SlotStatus = 1
	Chosen    SlotStatus = 2
)

type Slot struct {
	status SlotStatus
	value  string
}

func NewProposer(id, quorum int, value string, acceptors []string, generator IDService) *Proposer {
	proposer := &Proposer{
		id:                id,
		quorum:            quorum,
		acceptors:         make(map[string]pb.PaxosClient),
		proposalGenerator: generator,
	}
	for _, acceptor := range acceptors {
		cc, _ := grpc.Dial(acceptor,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.DefaultConfig,
			}),
		)
		proposer.acceptors[acceptor] = pb.NewPaxosClient(cc)
	}
	return proposer
}

func (p *Proposer) String() string {
	b := strings.Builder{}
	for _, s := range p.log {
		b.WriteString(s.value)
		b.WriteByte(' ')
	}
	return b.String()
}

func (p *Proposer) Start(in <-chan string) {
	defer p.wg.Wait()
	for {
		value, ok := <-in
		if !ok {
			log.Printf("[P%d] stopped, unable to read from input channel", p.id)
			return
		}

		// 一直循环直到当前值被选定
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				chosen := p.chooseOne(value)
				if value == chosen {
					break
				}
			}
		}()
	}
}

func (p *Proposer) returnFirstUnChosen() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := p.largestChosenIndex + 1; i < int64(len(p.log)); i++ {
		if p.log[i].status == Unused {
			return int64(i)
		}
	}

	slot := new(Slot)
	slot.status = Chooseing
	p.log = append(p.log, slot)
	return int64(len(p.log) - 1)
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
			p.log[index].status = Chosen
			p.log[index].value = value
			if index > p.largestChosenIndex {
				p.largestChosenIndex = index
			}
			log.Printf("[P%d] value at index %v is chosen: %v ", p.id, index, value)
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
		go func(addr string, acceptor pb.PaxosClient) {
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
			log.Printf("[P%d] Prepare on acceptor %v failed %v", p.id, addr, err)
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
			log.Printf("[P%d] Accept on acceptor %v failed %v", p.id, addr, err)
			return nil
		}
		return reply
	}
	return nil
}
