package store

import (
	pb "github.com/fudute/paxos/protoc"
)

type LogEntry struct {
	ID               int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	MiniProposal     int64  `gorm:"column:mini_proposal" json:"mini_proposal"`
	AcceptedProposal int64  `gorm:"column:accepted_proposal" json:"accepted_proposal"`
	AcceptedValue    string `gorm:"column:accepted_value" json:"accepted_value"`
	IsChosen         bool   `gorm:"column:is_chosen"`
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
