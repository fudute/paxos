package paxos

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
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

	// leader 信息，在leader信息有效时间内，只会接受来自leader的prepare和accept请求
	leader   string
	expire   *time.Timer
	interval time.Duration

	noMoreAccepteded    sync.Map
	noMoreAcceptedCount int32

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
		interval:          time.Second * 3,
		config:            config,
		conns:             NewConnectionsManager(),
		noMoreAccepteded:  sync.Map{},
		store:             store,
		logger:            log.New(os.Stderr, fmt.Sprintf("[%v]", name), 0),
	}
	return proposer
}

func (s *PaxosService) LeaderExpired() bool {
	if s.leader == "" {
		return true
	}
	select {
	case <-s.expire.C:
		return true
	default:
		return false
	}
}

func (s *PaxosService) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.ProposeReply, error) {
	if s.config.Leader && s.leader != s.self && !s.LeaderExpired() {
		cli, err := s.conns.Proposer(s.leader)
		if err == nil {
			return cli.Propose(ctx, req)
		}
	}
	for {
		chosen, err := s.chooseOne(req.ProposalValue)
		if err != nil {
			s.logger.Print("choose filed ", err)
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
	// s.logger.Print("start to chose at ", index)
	var proposalNum int64
	for {
		// s.logger.Println("use propose number ", proposalNum)

		if int(s.noMoreAcceptedCount) < s.quorum {
			proposalNum = s.proposalGenerator.Next()
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
				if reply.NoMoreAccepted {
					_, loaded := s.noMoreAccepteded.LoadOrStore(reply.From, nil)
					if loaded {
						atomic.AddInt32(&s.noMoreAcceptedCount, 1)
					}
				}
			}
		}

		acceptReplys := s.broadcast(s.config.Cluster.Nodes,
			func(node *config.Node) interface{} {
				return s.onAccept(index, proposalNum, value, node)
			})

		isChosen := true
		for i := 0; i <= s.quorum; i++ {
			reply := (<-acceptReplys).(*pb.AcceptReply)
			if reply.GetMiniProposal() > proposalNum {
				isChosen = false
				_, loaded := s.noMoreAccepteded.LoadAndDelete(reply.From)
				if loaded {
					atomic.AddInt32(&s.noMoreAcceptedCount, -1)
				}
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
	reply, err := cli.Prepare(context.Background(), &pb.PrepareRequest{
		Index:       index,
		ProposalNum: proposalNum,
	})
	if err != nil {
		s.logger.Printf("Prepare on acceptor %v failed %v", node.Addr, err)
		return nil
	}
	// s.logger.Printf("prepare on %v acceptedProposal %v acceptedValue %v", node.Name, reply.AcceptedProposal, reply.AcceptedValue)
	return reply
}

func (s *PaxosService) onAccept(index, proposalNum int64, value string, node *config.Node) *pb.AcceptReply {
	cli, err := s.conns.Acceptor(node.Addr)
	if err != nil {
		return nil
	}
	reply, err := cli.Accept(context.Background(), &pb.AcceptRequest{
		Index:         index,
		ProposalNum:   proposalNum,
		ProposalValue: value,
	})
	if err != nil {
		s.logger.Printf("Accept on acceptor %v failed %v", node.Addr, err)
		return nil
	}
	// s.logger.Printf("accept on %v, reply %v", node.Name, reply.MiniProposal)
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

	if s.config.Leader {
		if s.leader != req.Proposer {
			log.Printf("[P%v] leader change from %v to %v", s.self, s.leader, req.Proposer)
		}

		s.leader = req.Proposer
		if s.expire == nil {
			s.expire = time.NewTimer(s.interval)
		} else {
			s.expire.Reset(s.interval)
		}
	}

	s.logger.Printf("learn %v %v", req.Index, req.ProposalValue)
	return s.store.Learn(req)
}
