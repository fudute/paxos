package main

import (
	"fmt"
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
		log.Fatal("unable to load config")
	}

	var portStart = 8888
	var nProposer, nAcceptor = 2, 5

	var quorum = (nAcceptor + 1) / 2

	var acceptorsAddr, proposersAddr []string
	for i := 0; i < nAcceptor; i++ {
		acceptorsAddr = append(acceptorsAddr, fmt.Sprint(":", portStart+i))
	}
	for i := 0; i < nProposer; i++ {
		proposersAddr = append(proposersAddr, fmt.Sprint(":", portStart+nAcceptor+i))
	}

	var wg sync.WaitGroup
	for _, acceptor := range conf.Cluster.Acceptors {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			s := grpc.NewServer()
			s.RegisterService(&pb.Acceptor_ServiceDesc, paxos.NewAcceptor())
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("acceptor listen at ", addr)
			log.Fatal(s.Serve(lis))
		}(acceptor.Addr)
	}

	idService := NewUUIDService()

	var proposers []*paxos.Proposer

	for _, proposer := range conf.Cluster.Proposers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			s := grpc.NewServer()
			proposer := paxos.NewProposer(addr, quorum, &conf.Cluster, idService)
			proposers = append(proposers, proposer)
			s.RegisterService(&pb.Proposer_ServiceDesc, proposer)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("proposer listen at ", addr)
			log.Fatal(s.Serve(lis))
		}(proposer.Addr)
	}

	wg.Wait()
}
