package paxos

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	grpc "google.golang.org/grpc"
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
	acceptors         map[string]PaxosClient
}

func NewProposer(id, quorum int, value string, acceptors []string, generator IDService) *Proposer {
	proposer := &Proposer{
		id:                id,
		quorum:            quorum,
		value:             value,
		acceptors:         make(map[string]PaxosClient),
		proposalGenerator: generator,
	}
	for _, acceptor := range acceptors {
		cc, err := grpc.Dial(acceptor,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatal("unable to connect to acceptor ", acceptor, err)
		}
		proposer.acceptors[acceptor] = NewPaxosClient(cc)
	}
	return proposer
}

func (p *Proposer) Start() {
	log := log.New(os.Stdout, fmt.Sprintf("[P%v] ", p.id), 0)
	for {
		n := p.proposalGenerator.Next()
		prepareReplys := p.broadcast(n, p.onPrepare)

		var mini int64
		for i := 0; i <= p.quorum; i++ {
			reply := (<-prepareReplys).(*PrepareReply)
			if reply.AcceptedProposal > mini && reply.AcceptedValue != "" {
				p.value = reply.AcceptedValue
				mini = reply.AcceptedProposal
			}
		}

		acceptReplys := p.broadcast(n, p.onAccept)

		isChosen := true
		for i := 0; i <= p.quorum; i++ {
			reply := (<-acceptReplys).(*AcceptReply)
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
		go func(addr string, acceptor PaxosClient) {
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
func (p *Proposer) onPrepare(n int64, addr string) interface{} {
	reply, err := p.acceptors[addr].Prepare(context.Background(), &PrepareRequest{ProposalNum: n})
	if err != nil {
		log.Printf("[P%d] Prepare on acceptor %v failed %v", p.id, addr, err)
		return nil
	}
	return reply
}

func (p *Proposer) onAccept(n int64, addr string) interface{} {
	reply, err := p.acceptors[addr].Accept(context.Background(), &AcceptRequest{
		ProposalNum:   n,
		ProposalValue: p.value,
	})
	if err != nil {
		log.Printf("[P%d] Accept on acceptor %v failed %v", p.id, addr, err)
		return nil
	}
	return reply
}
