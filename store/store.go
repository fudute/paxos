package store

import (
	pb "github.com/fudute/paxos/protoc"
)

type LogEntry struct {
	ID               int64
	MiniProposal     int64
	AcceptedProposal int64
	AcceptedValue    string
	IsChosen         bool
}

func (*LogEntry) TableName() string {
	return "paxos_log"
}

// 日志
type LogStore interface {
	Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error)
	Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error)
	Learn(req *pb.LearnRequest) (*pb.LearnReply, error)
	PickSlot() int64
}
