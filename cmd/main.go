package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fudute/paxos"
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
	flag.Parse()

	var portStart = 8888
	var nProposer, nAcceptor = 2, 5

	var quorum = (nAcceptor + 1) / 2

	var peers []string
	for i := 0; i < nAcceptor; i++ {
		peers = append(peers, fmt.Sprint(":", portStart+i))
	}

	for _, peer := range peers {
		go func(peer string) {
			s := grpc.NewServer()
			s.RegisterService(&pb.Paxos_ServiceDesc, paxos.NewAcceptor())
			lis, err := net.Listen("tcp", peer)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("server listen at ", peer)
			log.Fatal(s.Serve(lis))
		}(peer)
	}

	idService := NewUUIDService()

	in := make(chan string, 3)
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			<-time.After(time.Second / 100)
			in <- fmt.Sprint(i)
		}
	}()

	var wg sync.WaitGroup
	var proposers []*paxos.Proposer
	for i := 0; i < nProposer; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			proposer := paxos.NewProposer(i, quorum, fmt.Sprint(i), peers, idService)
			proposers = append(proposers, proposer)
			proposer.Start(in)
		}(i)
	}
	wg.Wait()

	for _, p := range proposers {
		fmt.Println(p.String())
	}
}
