package main

import (
	"context"
	"net"

	"github.com/fudute/paxos/paxos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type paxosServer struct {
	paxos.UnimplementedPaxosServer
}

func (*paxosServer) Prepare(context.Context, *paxos.PrepareRequest) (*paxos.PrepareReply, error) {
	return &paxos.PrepareReply{}, nil
}
func (*paxosServer) Accept(context.Context, *paxos.AcceptRequest) (*paxos.AcceptReply, error) {
	return &paxos.AcceptReply{}, nil
}

func main() {
	s := grpc.NewServer()
	s.RegisterService(&paxos.Paxos_ServiceDesc, &paxosServer{})

	lis, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		grpclog.Fatal("listen failed", err)
	}

	grpclog.Info("server listen at 127.0.0.1:8888")
	grpclog.Error(s.Serve(lis))
}
