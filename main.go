package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/fudute/paxos/paxos"
	"google.golang.org/grpc"
)

var peers = []string{":8888", ":8889", ":8890"}

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
	flag.Parse()

	for _, peer := range peers {
		go func(peer string) {
			s := grpc.NewServer()
			s.RegisterService(&paxos.Paxos_ServiceDesc, paxos.NewService())
			lis, err := net.Listen("tcp", peer)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("server listen at ", peer)
			log.Fatal(s.Serve(lis))
		}(peer)
	}

	idService := NewUUIDService()

	var wg sync.WaitGroup
	var n int = 3
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			proposer := paxos.NewProposer(i, (n+1)/2, fmt.Sprint(i), peers, idService)
			proposer.Start()
		}(i)
	}
	wg.Wait()
}
