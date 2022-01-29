package store

import (
	"bytes"
	"encoding/gob"
	"sync"
	"sync/atomic"

	pb "github.com/fudute/paxos/protoc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type levelDBLogStore struct {
	mu sync.Mutex
	db *leveldb.DB

	buf             *bytes.Buffer
	largestAccepted int64
	largestChosen   int64

	wo *opt.WriteOptions
	ro *opt.ReadOptions
}

func (l *levelDBLogStore) get(key interface{}, value interface{}) error {
	defer l.buf.Reset()

	enc := gob.NewEncoder(l.buf)
	err := enc.Encode(key)
	if err != nil {
		return err
	}

	bs, err := l.db.Get(l.buf.Bytes(), l.ro)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(bs)).Decode(value)
}

func (l *levelDBLogStore) set(key interface{}, value interface{}) error {
	defer l.buf.Reset()

	enc := gob.NewEncoder(l.buf)
	err := enc.Encode(key)
	if err != nil {
		return err
	}

	ind := l.buf.Len()
	err = enc.Encode(value)
	if err != nil {
		return err
	}
	return l.db.Put(l.buf.Bytes()[:ind], l.buf.Bytes()[ind+1:], l.wo)
}

func (l *levelDBLogStore) Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ins := new(pb.LogEntry)

	err := l.get(req.GetIndex(), ins)
	if err != nil {
		return nil, err
	}

	if req.GetProposalNum() > ins.MiniProposal {
		ins.MiniProposal = req.GetProposalNum()
	}

	l.set(req.GetIndex(), ins)

	return &pb.PrepareReply{
		AcceptedProposal: ins.AcceptedProposal,
		AcceptedValue:    ins.AcceptedValue,
		NoMoreAccepted:   l.largestAccepted < req.Index,
	}, nil
}

func (l *levelDBLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ins := new(pb.LogEntry)
	err := l.get(req.GetIndex(), ins)
	if err != nil {
		return nil, err
	}

	if req.GetProposalNum() >= ins.MiniProposal {
		ins.MiniProposal = req.GetProposalNum()
		ins.AcceptedProposal = ins.MiniProposal
		ins.AcceptedValue = req.GetProposalValue()
		if l.largestAccepted < req.Index {
			l.largestAccepted = req.Index
		}
	}

	err = l.set(req.GetIndex(), ins)
	if err != nil {
		return nil, err
	}

	return &pb.AcceptReply{
		MiniProposal: ins.MiniProposal,
	}, nil

}
func (l *levelDBLogStore) Learn(req *pb.LearnRequest) (*pb.LearnReply, error) {
	err := l.set(req.GetIndex(), &LogEntry{
		AcceptedValue: req.GetProposalValue(),
		IsChosen:      true,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LearnReply{}, nil
}

func (l *levelDBLogStore) PickSlot() int64 {
	return atomic.AddInt64(&l.largestChosen, 1)
}
