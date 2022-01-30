package store

import (
	"sync/atomic"

	"github.com/cenk/backoff"
	pb "github.com/fudute/paxos/protoc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlLogStore struct {
	db              *gorm.DB
	largestAccepted int64
	picker          Picker
}

func NewMySQLLogStore(dsn string) (*MysqlLogStore, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	db.Migrator().DropTable(LogEntry{})
	db.AutoMigrate(LogEntry{})

	return &MysqlLogStore{
		db: db,
	}, nil
}

func (l *MysqlLogStore) Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	ins := LogEntry{}

	err := backoff.Retry(func() error {
		return l.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).FirstOrCreate(&ins, LogEntry{ID: req.Index}).Error; err != nil {
				return err
			}
			if req.GetProposalNum() > ins.MiniProposal {
				ins.MiniProposal = req.GetProposalNum()
				return tx.Model(LogEntry{}).Where("id=?", req.Index).Updates(LogEntry{MiniProposal: req.GetProposalNum()}).Error
			}
			return nil
		})
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return nil, err
	}
	return &pb.PrepareReply{
		AcceptedProposal: ins.AcceptedProposal,
		AcceptedValue:    ins.AcceptedValue,
		NoMoreAccepted:   l.largestAccepted < req.Index,
	}, nil
}

func (l *MysqlLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	ins := LogEntry{}

	err := backoff.Retry(func() error {
		return l.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).FirstOrCreate(&ins, LogEntry{ID: req.Index}).Error; err != nil {
				return err
			}
			if req.GetProposalNum() >= ins.MiniProposal {
				ins.MiniProposal = req.GetProposalNum()
				ins.AcceptedProposal = ins.MiniProposal
				ins.AcceptedValue = req.GetProposalValue()

				for l.largestAccepted < req.Index {
					atomic.CompareAndSwapInt64(&l.largestAccepted, l.largestAccepted, req.Index)
				}
			}
			return tx.Where(LogEntry{ID: req.Index}).Updates(&ins).Error
		})
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return nil, err
	}
	return &pb.AcceptReply{
		MiniProposal: ins.MiniProposal,
	}, nil
}
func (l *MysqlLogStore) Learn(req *pb.LearnRequest) (*pb.LearnReply, error) {
	ins := LogEntry{
		ID:            req.Index,
		AcceptedValue: req.ProposalValue,
		IsChosen:      true,
	}

	err := l.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},                                     // key colume
		DoUpdates: clause.AssignmentColumns([]string{"accepted_value", "is_chosen"}), // column needed to be updated
	}).Create(&ins).Error

	l.picker.Put(int(req.GetIndex()))

	return &pb.LearnReply{}, err
}

func (l *MysqlLogStore) PickSlot() int64 {
	return int64(l.picker.Pick())
}
