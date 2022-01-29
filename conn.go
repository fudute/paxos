package paxos

import (
	"sync"

	pb "github.com/fudute/paxos/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionsManager struct {
	m sync.Map
}

func NewConnectionsManager() *ConnectionsManager {
	return &ConnectionsManager{}
}

func (cp *ConnectionsManager) GetConn(addr string) (grpc.ClientConnInterface, error) {
	if cc, ok := cp.m.Load(addr); ok {
		return cc.(grpc.ClientConnInterface), nil
	}

	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	actual, _ := cp.m.LoadOrStore(addr, cc)
	return actual.(grpc.ClientConnInterface), nil
}

func (cp *ConnectionsManager) Proposer(addr string) (pb.ProposerClient, error) {
	cc, err := cp.GetConn(addr)
	if err != nil {
		return nil, err
	}
	cli := pb.NewProposerClient(cc)
	return cli, nil
}
func (cp *ConnectionsManager) Acceptor(addr string) (pb.AcceptorClient, error) {
	cc, err := cp.GetConn(addr)
	if err != nil {
		return nil, err
	}
	cli := pb.NewAcceptorClient(cc)
	return cli, nil
}
func (cp *ConnectionsManager) Learner(addr string) (pb.LearnerClient, error) {
	cc, err := cp.GetConn(addr)
	if err != nil {
		return nil, err
	}
	cli := pb.NewLearnerClient(cc)
	return cli, nil
}
