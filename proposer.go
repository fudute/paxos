package paxos

import (
	"context"
	"fmt"
	"log"
	"os"
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
	value             string
	acceptors         map[string]pb.PaxosClient
}

func NewProposer(id, quorum int, value string, acceptors []string, generator IDService) *Proposer {
	proposer := &Proposer{
		id:                id,
		quorum:            quorum,
		value:             value,
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

func (p *Proposer) Start() {
	log := log.New(os.Stdout, fmt.Sprintf("[P%v] ", p.id), 0)
	for {
		n := p.proposalGenerator.Next()
		prepareReplys := p.broadcast(n, func(n int64, addr string) interface{} { return p.onPrepare(n, addr) })

		var mini int64
		for i := 0; i <= p.quorum; i++ {
			reply := (<-prepareReplys).(*pb.PrepareReply)
			if reply.AcceptedProposal > mini && reply.AcceptedValue != "" {
				p.value = reply.AcceptedValue
				mini = reply.AcceptedProposal
			}
		}

		acceptReplys := p.broadcast(n, func(n int64, addr string) interface{} { return p.onAccept(n, addr) })

		isChosen := true
		for i := 0; i <= p.quorum; i++ {
			reply := (<-acceptReplys).(*pb.AcceptReply)
			if reply.GetMiniProposal() > n {
				isChosen = false
				break
			}
		}

		if isChosen {
			log.Printf("value is chosen: %v ", p.value)
			return
		}
	}
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

func (p *Proposer) onPrepare(n int64, addr string) *pb.PrepareReply {
	if cli, ok := p.acceptors[addr]; ok {
		reply, err := cli.Prepare(context.Background(), &pb.PrepareRequest{ProposalNum: n})
		if err != nil {
			log.Printf("[P%d] Prepare on acceptor %v failed %v", p.id, addr, err)
			return nil
		}
		return reply
	}
	return nil
}

func (p *Proposer) onAccept(n int64, addr string) *pb.AcceptReply {
	if cli, ok := p.acceptors[addr]; ok {
		reply, err := cli.Accept(context.Background(), &pb.AcceptRequest{
			ProposalNum:   n,
			ProposalValue: p.value,
		})
		if err != nil {
			log.Printf("[P%d] Accept on acceptor %v failed %v", p.id, addr, err)
			return nil
		}
		return reply
	}
	return nil
}
