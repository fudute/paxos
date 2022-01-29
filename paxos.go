package paxos

import (
	"context"
	"fmt"
	"log"
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
	cluster           *config.Cluster
	conns             *ConnectionsManager

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

func NewPaxosService(name, addr string, quorum int, cluster *config.Cluster, generator IDService) *PaxosService {
	dsn := fmt.Sprintf("root:123@tcp(127.0.0.1:3306)/demo_%s?charset=utf8mb4&parseTime=True&loc=Local", name)
	store, err := store.NewMySQLLogStore(dsn)
	if err != nil {
		log.Fatal(err)
	}
	proposer := &PaxosService{
		name:              name,
		self:              addr,
		quorum:            quorum,
		proposalGenerator: generator,
		interval:          time.Second,
		cluster:           cluster,
		conns:             NewConnectionsManager(),
		noMoreAccepteded:  sync.Map{},
		store:             store,
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
	if s.leader != s.self && !s.LeaderExpired() {
		cli, err := s.conns.Proposer(s.leader)
		if err == nil {
			return cli.Propose(ctx, req)
		}
	}
	for {
		chosen := s.chooseOne(req.ProposalValue)
		if req.ProposalValue == chosen {
			break
		}
	}
	return &pb.ProposeReply{}, nil
}

func (s *PaxosService) chooseOne(value string) (chosen string) {
	index := s.store.PickSlot()
	for {
		proposalNum := s.proposalGenerator.Next()

		if int(s.noMoreAcceptedCount) < s.quorum {
			prepareReplys := s.broadcast(s.cluster.Nodes,
				func(node *config.Node) interface{} {
					return s.onPrepare(index, proposalNum, node)
				})

			var mini int64
			for i := 0; i <= s.quorum; i++ {
				reply, ok := (<-prepareReplys).(*pb.PrepareReply)
				// 如果失败了，等待一段时间后重试
				if !ok {
					<-time.After(time.Second)
					continue
				}
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

		acceptReplys := s.broadcast(s.cluster.Nodes,
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
			s.broadcast(s.cluster.Nodes, func(node *config.Node) interface{} {
				return s.OnLearn(index, value, node)
			})
			break
		}
	}
	return value
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
		log.Printf("[P%s] Prepare on acceptor %v failed %v", s.self, node.Addr, err)
		return nil
	}
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
		log.Printf("[P%s] Accept on acceptor %v failed %v", s.self, node.Addr, err)
		return nil
	}
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
		log.Printf("[P%s] Learn on learner %v failed %v", s.self, node.Addr, err)
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
	log.Printf("[%v] value at %v is chosen %v", s.name, req.Index, req.ProposalValue)
	return s.store.Learn(req)
}
