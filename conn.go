package paxos

import (
	"sync"

	pb "github.com/fudute/paxos/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionsManager struct {
	m  map[string]grpc.ClientConnInterface
	mu sync.Mutex
}

func NewConnectionsManager() *ConnectionsManager {
	return &ConnectionsManager{
		m: make(map[string]grpc.ClientConnInterface),
	}
}

func (cp *ConnectionsManager) GetConn(addr string) (grpc.ClientConnInterface, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cc, ok := cp.m[addr]; ok {
		return cc, nil
	}

	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cp.m[addr] = cc
	return cc, nil
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
