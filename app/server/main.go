package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

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

	var acceptorsAddr, proposersAddr []string
	for i := 0; i < nAcceptor; i++ {
		acceptorsAddr = append(acceptorsAddr, fmt.Sprint(":", portStart+i))
	}
	for i := 0; i < nProposer; i++ {
		proposersAddr = append(proposersAddr, fmt.Sprint(":", portStart+nAcceptor+i))
	}

	var wg sync.WaitGroup
	for _, acceptor := range acceptorsAddr {
		wg.Add(1)
		go func(acceptor string) {
			defer wg.Done()
			s := grpc.NewServer()
			s.RegisterService(&pb.Acceptor_ServiceDesc, paxos.NewAcceptor())
			lis, err := net.Listen("tcp", acceptor)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("acceptor listen at ", acceptor)
			log.Fatal(s.Serve(lis))
		}(acceptor)
	}

	idService := NewUUIDService()

	var proposers []*paxos.Proposer

	for _, proposerAddr := range proposersAddr {
		wg.Add(1)
		go func(proposerAddr string) {
			defer wg.Done()
			s := grpc.NewServer()
			proposer := paxos.NewProposer(proposerAddr, quorum, acceptorsAddr, proposersAddr, idService)
			proposers = append(proposers, proposer)
			s.RegisterService(&pb.Proposer_ServiceDesc, proposer)
			lis, err := net.Listen("tcp", proposerAddr)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("proposer listen at ", proposerAddr)
			log.Fatal(s.Serve(lis))
		}(proposerAddr)
	}

	// in := make(chan string, 3)
	// go func() {
	// 	defer close(in)
	// 	for i := 0; i < 10; i++ {
	// 		<-time.After(time.Second / 100)
	// 		in <- fmt.Sprint(i)
	// 	}
	// }()

	// var wg sync.WaitGroup
	// for _, p := range proposers {
	// 	wg.Add(1)
	// 	go func(p *paxos.Proposer) {
	// 		defer wg.Done()
	// 		p.Start(in)
	// 	}(p)
	// }
	// wg.Wait()

	// for _, p := range proposers {
	// 	fmt.Println(p.String())
	// }
	wg.Wait()
}
