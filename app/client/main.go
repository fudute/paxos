package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fudute/paxos/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var N = 10000

var wg sync.WaitGroup

func main() {
	start := time.Now()

	cc1, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cc2, err := grpc.Dial(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan string)

	go func() {
		for i := 0; i < N; i++ {
			ch <- fmt.Sprint(i)
		}
		close(ch)
	}()

	wg.Add(2)
	go send(protoc.NewProposerClient(cc1), ch)
	go send(protoc.NewProposerClient(cc2), ch)

	wg.Wait()

	fmt.Println("qps: ", float64(N)/time.Since(start).Seconds())
}

func send(cli protoc.ProposerClient, in <-chan string) {
	defer wg.Done()
	var err error
	for {
		value, ok := <-in
		if !ok {
			return
		}
		_, err = cli.Propose(context.Background(), &protoc.ProposeRequest{
			ProposalValue: value,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
