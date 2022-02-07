package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
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

	buf           *bytes.Buffer
	largestChosen int64

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
	log.Printf("[store] get %v %+v", index, entry)
	return entry, nil
}

func (l *LevelDBLogStore) set(index int64, value *pb.LogEntry) error {
	log.Printf("[store] set %v %+v", index, value)
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

	if ins.AcceptedProposal == 0 && ins.AcceptedValue != "" {
		log.Fatalf("bad at %v %+v", req.Index, ins)
	}

	return &pb.PrepareReply{
		AcceptedProposal: ins.AcceptedProposal,
		AcceptedValue:    ins.AcceptedValue,
	}, nil
}

func (l *LevelDBLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ins, err := l.get(req.GetIndex())
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return nil, err
		}
	}

	if req.GetProposalNum() >= ins.MiniProposal {
		if ins.Status == pb.LogEntryStatus_Chosen && req.ProposalValue != ins.AcceptedValue {
			log.Fatalf("chosen different value at %v, new: %v, old: %v", req.Index, req.ProposalValue, ins.AcceptedValue)
		}
		ins.MiniProposal = req.GetProposalNum()
		ins.AcceptedProposal = ins.MiniProposal
		ins.AcceptedValue = req.GetProposalValue()
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

	ins, err := l.get(req.GetIndex())
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return nil, err
		}
	}
	if req.GetProposalNum() >= ins.MiniProposal {
		err := l.set(req.GetIndex(), &pb.LogEntry{
			MiniProposal:     req.GetProposalNum(),
			AcceptedProposal: req.GetProposalNum(),
			AcceptedValue:    req.GetProposalValue(),
			Status:           pb.LogEntryStatus_Chosen,
		})
		if err != nil {
			return nil, err
		}
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
