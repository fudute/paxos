package main

import (
	"context"

	"github.com/fudute/paxos/paxos"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("127.0.0.1:8888", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		glog.Fatal("grpc dial err", err)
	}
	cli := paxos.NewPaxosClient(cc)

	reply, err := cli.Accept(context.Background(), &paxos.AcceptRequest{})
	if err != nil {
		glog.Error(err)
	}
	glog.Info("reply ", reply)
}
