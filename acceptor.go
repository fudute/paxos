package paxos

import (
	context "context"
	"sync"

	pb "github.com/fudute/paxos/protoc"
)

type acceptor struct {
	mu sync.Mutex
	pb.UnimplementedPaxosServer
	miniProposal     int64
	acceptedProposal int64
	acceptedValue    string
}

func NewAcceptor() pb.PaxosServer {
	return &acceptor{}
}

func (a *acceptor) Prepare(ctx context.Context, req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if req.GetProposalNum() > a.miniProposal {
		a.miniProposal = req.GetProposalNum()
	}

	return &pb.PrepareReply{
		AcceptedProposal: a.acceptedProposal,
		AcceptedValue:    a.acceptedValue,
	}, nil
}
func (a *acceptor) Accept(ctx context.Context, req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if req.GetProposalNum() >= a.miniProposal {
		a.miniProposal = req.GetProposalNum()
		a.acceptedProposal = a.miniProposal
		a.acceptedValue = req.GetProposalValue()
	}

	return &pb.AcceptReply{
		MiniProposal: a.miniProposal,
	}, nil
}
