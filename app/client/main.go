package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fudute/paxos/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc1, err := grpc.Dial(":8893", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cc2, err := grpc.Dial(":8894", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan string)

	go func() {
		for i := 0; i < 100; i++ {
			ch <- fmt.Sprint(i)
		}
		close(ch)
	}()

	go send(protoc.NewProposerClient(cc1), ch)
	go send(protoc.NewProposerClient(cc2), ch)

	fmt.Scan(new(byte))
}

func send(cli protoc.ProposerClient, in <-chan string) {
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
