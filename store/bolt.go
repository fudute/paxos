package store

import (
	"strconv"

	"github.com/boltdb/bolt"
	pb "github.com/fudute/paxos/protoc"
	"google.golang.org/protobuf/proto"
)

type boltLogStore struct {
	db         *bolt.DB
	bucketName []byte
}

func NewBoltLogStore(path string) (LogStore, error) {
	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists([]byte("paxos_log"))
		return err
	})
	if err != nil {
		return nil, err
	}
	return &boltLogStore{
		db:         db,
		bucketName: []byte("paxos"),
	}, nil
}

func instanceByKey(b *bolt.Bucket, key []byte) (*pb.LogEntry, error) {
	value := b.Get(key)
	ins := new(pb.LogEntry)
	err := proto.Unmarshal(value, ins)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (b *boltLogStore) Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	reply := new(pb.PrepareReply)
	err := b.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(b.bucketName)

		key := []byte(strconv.Itoa(int(req.GetIndex())))

		ins, err := instanceByKey(b, key)
		if err != nil {
			return err
		}

		if req.GetProposalNum() > ins.MiniProposal {
			ins.MiniProposal = req.GetProposalNum()
		}
		reply.AcceptedProposal = ins.AcceptedProposal
		reply.AcceptedValue = ins.AcceptedValue

		value, _ := proto.Marshal(ins)
		return b.Put(key, value)
	})

	return reply, err
}

func (b *boltLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	reply := new(pb.AcceptReply)

	return reply, nil

}
func (b *boltLogStore) Learn(req *pb.LearnRequest) (*pb.LearnReply, error) {
	reply := new(pb.LearnReply)

	return reply, nil
}
func (b *boltLogStore) PickSlot() int64 {
	return 0
}
