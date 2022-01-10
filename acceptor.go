package paxos

import (
	context "context"
	"sync"

	pb "github.com/fudute/paxos/protoc"
)

type acceptor struct {
	pb.UnimplementedPaxosServer

	ins sync.Map // key is index (int64), value is *Instance
}

type Instance struct {
	mu               sync.Mutex
	miniProposal     int64
	acceptedProposal int64
	acceptedValue    string
}

func NewAcceptor() pb.PaxosServer {
	return &acceptor{}
}

func (a *acceptor) Prepare(ctx context.Context, req *pb.PrepareRequest) (*pb.PrepareReply, error) {

	actual, _ := a.ins.LoadOrStore(req.Index, new(Instance))
	ins := actual.(*Instance)

	ins.mu.Lock()
	defer ins.mu.Unlock()

	if req.GetProposalNum() > ins.miniProposal {
		ins.miniProposal = req.GetProposalNum()
	}

	return &pb.PrepareReply{
		AcceptedProposal: ins.acceptedProposal,
		AcceptedValue:    ins.acceptedValue,
	}, nil
}

func (a *acceptor) Accept(ctx context.Context, req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	actual, _ := a.ins.LoadOrStore(req.Index, new(Instance))
	ins := actual.(*Instance)

	ins.mu.Lock()
	defer ins.mu.Unlock()

	if req.GetProposalNum() >= ins.miniProposal {
		ins.miniProposal = req.GetProposalNum()
		ins.acceptedProposal = ins.miniProposal
		ins.acceptedValue = req.GetProposalValue()
	}

	return &pb.AcceptReply{
		MiniProposal: ins.miniProposal,
	}, nil
}
