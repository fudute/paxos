package paxos

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fudute/paxos/config"
	pb "github.com/fudute/paxos/protoc"
	"github.com/fudute/paxos/store"
)

type IDService interface {
	Next() int64
}

type PaxosService struct {
	name              string
	self              string
	quorum            int
	proposalGenerator IDService
	config            *config.Config
	conns             *ConnectionsManager
	logger            *log.Logger

	store store.LogStore

	pb.UnimplementedProposerServer
	pb.UnimplementedAcceptorServer
	pb.UnimplementedLearnerServer
}

func NewPaxosService(name, addr string, quorum int, config *config.Config, generator IDService, store store.LogStore) *PaxosService {
	proposer := &PaxosService{
		name:              name,
		self:              addr,
		quorum:            quorum,
		proposalGenerator: generator,
		config:            config,
		conns:             NewConnectionsManager(),
		store:             store,
		logger:            log.New(os.Stderr, fmt.Sprintf("[%v]", name), 0),
	}
	return proposer
}

func (s *PaxosService) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeReply, error) {
	for {
		chosen, err := s.chooseOne(req.ProposalValue)
		if err != nil {
			s.logger.Print("choose failed ", err)
			<-time.After(time.Second)
			continue
		}
		if req.ProposalValue == chosen {
			break
		}
	}
	return &pb.ProposeReply{}, nil
}

func (s *PaxosService) chooseOne(value string) (chosen string, err error) {
	index := s.store.PickSlot()
	s.logger.Print("start to chose at ", index)
	for {
		proposalNum := s.proposalGenerator.Next()
		s.logger.Println("use propose number ", proposalNum)

		s.logger.Printf("prepare at %v with proposal %v", index, proposalNum)
		prepareReplys := s.broadcast(s.config.Cluster.Nodes,
			func(node *config.Node) interface{} {
				return s.onPrepare(index, proposalNum, node)
			})

		var mini int64
		for i := 0; i <= s.quorum; i++ {
			reply := (<-prepareReplys).(*pb.PrepareReply)
			if reply.AcceptedProposal > mini && reply.AcceptedValue != "" {
				value = reply.AcceptedValue
				mini = reply.AcceptedProposal
			}
		}

		s.logger.Printf("accept at %v with proposal %v value %v", index, proposalNum, value)
		acceptReplys := s.broadcast(s.config.Cluster.Nodes,
			func(node *config.Node) interface{} {
				return s.onAccept(index, proposalNum, value, node)
			})

		isChosen := true
		for i := 0; i <= s.quorum; i++ {
			reply := (<-acceptReplys).(*pb.AcceptReply)
			if reply.GetMiniProposal() > proposalNum {
				isChosen = false
				break
			}
		}

		if isChosen {
			s.logger.Printf("value at %v is chosen %v", index, value)
			s.broadcast(s.config.Cluster.Nodes, func(node *config.Node) interface{} {
				return s.OnLearn(index, value, node)
			})
			break
		}
	}
	return value, nil
}

func (s *PaxosService) broadcast(nodes []*config.Node, f func(node *config.Node) interface{}) <-chan interface{} {
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

func (s *PaxosService) onPrepare(index, proposalNum int64, node *config.Node) *pb.PrepareReply {
	cli, err := s.conns.Acceptor(node.Addr)
	if err != nil {
		return nil
	}
	req := &pb.PrepareRequest{
		Index:       index,
		ProposalNum: proposalNum,
	}
	reply, err := cli.Prepare(context.Background(), req)
	if err != nil {
		s.logger.Printf("Prepare on acceptor %v failed %v", node.Addr, err)
		return nil
	}
	s.logger.Printf("prepare on %v req %+v reply %+v", node.Name, req, reply)
	return reply
}

func (s *PaxosService) onAccept(index, proposalNum int64, value string, node *config.Node) *pb.AcceptReply {
	cli, err := s.conns.Acceptor(node.Addr)
	if err != nil {
		return nil
	}

	req := &pb.AcceptRequest{
		Index:         index,
		ProposalNum:   proposalNum,
		ProposalValue: value,
	}

	reply, err := cli.Accept(context.Background(), req)
	if err != nil {
		s.logger.Printf("Accept on acceptor %v failed %v", node.Addr, err)
		return nil
	}
	s.logger.Printf("accept on %v req %+v reply %+v", node.Name, req, reply)
	return reply
}

func (s *PaxosService) OnLearn(index int64, value string, node *config.Node) *pb.LearnReply {
	cli, err := s.conns.Learner(node.Addr)
	if err != nil {
		return nil
	}
	reply, err := cli.Learn(context.Background(), &pb.LearnRequest{
		Index:         index,
		Proposer:      s.self,
		ProposalValue: value,
	})
	if err != nil {
		s.logger.Printf("Learn on learner %v failed %v", node.Addr, err)
		return nil
	}
	return reply
}

func (s *PaxosService) Prepare(ctx context.Context, req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	return s.store.Prepare(req)
}

func (s *PaxosService) Accept(ctx context.Context, req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	return s.store.Accept(req)
}

func (s *PaxosService) Learn(ctx context.Context, req *pb.LearnRequest) (*pb.LearnReply, error) {

	s.logger.Printf("learn %v %v", req.Index, req.ProposalValue)
	return s.store.Learn(req)
}
