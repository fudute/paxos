package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"

	pb "github.com/fudute/paxos/protoc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"google.golang.org/protobuf/proto"
)

type LevelDBLogStore struct {
	mu sync.Mutex
	db *leveldb.DB

	buf             *bytes.Buffer
	largestAccepted int64
	largestChosen   int64

	wo *opt.WriteOptions
	ro *opt.ReadOptions
}

func NewLevelDBLogStore(path string) (*LevelDBLogStore, error) {
	exec.Command("rm", "-rf", path).Run()
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDBLogStore{db: db, buf: bytes.NewBuffer(nil)}, nil
}

func (l *LevelDBLogStore) get(index int64) (*pb.LogEntry, error) {
	key := []byte(strconv.Itoa(int(index)))
	bs, err := l.db.Get(key, l.ro)
	if err != nil {
		return nil, err
	}
	entry := new(pb.LogEntry)
	if err := proto.Unmarshal(bs, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (l *LevelDBLogStore) set(index int64, value *pb.LogEntry) error {
	key := []byte(strconv.Itoa(int(index)))
	bs, err := proto.Marshal(value)
	if err != nil {
		return err
	}
	return l.db.Put(key, bs, l.wo)
}

func (l *LevelDBLogStore) Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ins, err := l.get(req.GetIndex())
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return nil, err
		}
		ins = new(pb.LogEntry)
	}

	if req.GetProposalNum() > ins.MiniProposal {
		ins.MiniProposal = req.GetProposalNum()
	}

	err = l.set(req.GetIndex(), ins)
	if err != nil {
		return nil, err
	}

	return &pb.PrepareReply{
		AcceptedProposal: ins.AcceptedProposal,
		AcceptedValue:    ins.AcceptedValue,
	}, nil
}

func (l *LevelDBLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ins := new(pb.LogEntry)
	ins, err := l.get(req.GetIndex())
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return nil, err
		}
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
func (l *LevelDBLogStore) Learn(req *pb.LearnRequest) (*pb.LearnReply, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	err := l.set(req.GetIndex(), &pb.LogEntry{
		AcceptedValue: req.GetProposalValue(),
		Status:        pb.LogEntryStatus_Chosen,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LearnReply{}, nil
}

func (l *LevelDBLogStore) PickSlot() int64 {
	return atomic.AddInt64(&l.largestChosen, 1)
}
func (l *LevelDBLogStore) Range() ([]*LogEntry, error) {

	sh, err := l.db.GetSnapshot()
	if err != nil {
		return nil, err
	}

	var res []*LogEntry
	iter := sh.NewIterator(nil, l.ro)
	for iter.Next() {
		var entry LogEntry
		gob.NewDecoder(bytes.NewReader(iter.Value())).Decode(&entry)
		res = append(res, &entry)
	}
	return res, nil
}
