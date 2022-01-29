package main

import (
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/fudute/paxos"
	"github.com/fudute/paxos/config"
	pb "github.com/fudute/paxos/protoc"
	"google.golang.org/grpc"
)

type UUID struct {
	id int64
}

func NewUUIDService() paxos.IDService {
	return &UUID{}
}

func (u *UUID) Next() int64 {
	return atomic.AddInt64(&u.id, 1)
}

func main() {
	log.SetFlags(0)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	log.Printf("load config: %#v", conf)
	var nAcceptor = len(conf.Cluster.Nodes)

	var quorum = (nAcceptor + 1) / 2

	var wg sync.WaitGroup
	// for _, acceptor := range conf.Cluster.Acceptors {
	// 	wg.Add(1)
	// 	go func(addr string) {
	// 		defer wg.Done()
	// 		s := grpc.NewServer()
	// 		s.RegisterService(&pb.Acceptor_ServiceDesc, paxos.NewAcceptor())
	// 		lis, err := net.Listen("tcp", addr)
	// 		if err != nil {
	// 			log.Fatal("listen failed ", err)
	// 		}
	// 		log.Println("acceptor listen at ", addr)
	// 		log.Fatal(s.Serve(lis))
	// 	}(acceptor.Addr)
	// }

	idService := NewUUIDService()

	for _, node := range conf.Cluster.Nodes {
		wg.Add(1)
		go func(node *config.Node) {
			defer wg.Done()
			s := grpc.NewServer()
			service := paxos.NewPaxosService(node.Name, node.Addr, quorum, conf, idService)
			s.RegisterService(&pb.Proposer_ServiceDesc, service)
			s.RegisterService(&pb.Learner_ServiceDesc, service)
			s.RegisterService(&pb.Acceptor_ServiceDesc, service)
			lis, err := net.Listen("tcp", node.Addr)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("proposer listen at ", node.Addr)
			log.Fatal(s.Serve(lis))
		}(node)
	}

	wg.Wait()
}
