package paxos

import (
	context "context"
	"sync"
)

type acceptor struct {
	mu sync.Mutex
	UnimplementedPaxosServer
	miniProposal     int64
	acceptedProposal int64
	acceptedValue    string
}

func NewService() PaxosServer {
	return &acceptor{}
}

func (a *acceptor) Prepare(ctx context.Context, req *PrepareRequest) (*PrepareReply, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if req.GetProposalNum() > a.miniProposal {
		a.miniProposal = req.GetProposalNum()
	}

	return &PrepareReply{
		AcceptedProposal: a.acceptedProposal,
		AcceptedValue:    a.acceptedValue,
	}, nil
}
func (a *acceptor) Accept(ctx context.Context, req *AcceptRequest) (*AcceptReply, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if req.GetProposalNum() >= a.miniProposal {
		a.miniProposal = req.GetProposalNum()
		a.acceptedProposal = a.miniProposal
		a.acceptedValue = req.GetProposalValue()
	}

	return &AcceptReply{
		MiniProposal: a.miniProposal,
	}, nil
}
