package store

import (
	"sync"
)

type MemLogStore struct {
	m  map[int64]LogEntry
	mu sync.Mutex
}

// func (ms *MemLogStore) Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error) {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()
// 	entry, ok := ms.m[req.Index]
// 	if ok {

// 	}
// }
// func (ms *MemLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {}
// func (ms *MemLogStore) Learn(req *pb.LearnRequest) (*pb.LearnReply, error)    {}
// func (ms *MemLogStore) Range() ([]*LogEntry, error)                           {}
// func (ms *MemLogStore) PickSlot() int64                                       {}
