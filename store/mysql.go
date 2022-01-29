package store

import (
	"sync/atomic"

	pb "github.com/fudute/paxos/protoc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlLogStore struct {
	db              *gorm.DB
	largestAccepted int64
	picker          Picker
}

func NewMySQLLogStore(dsn string) (LogStore, error) {
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		return nil, err
	}
	db.Migrator().DropTable(LogEntry{})
	db.AutoMigrate(LogEntry{})

	return &mysqlLogStore{
		db: db,
	}, nil
}

func (l *mysqlLogStore) Prepare(req *pb.PrepareRequest) (*pb.PrepareReply, error) {
	ins := LogEntry{}

	err := l.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.FirstOrCreate(&ins, LogEntry{ID: req.Index}).Error; err != nil {
			return err
		}
		if req.GetProposalNum() > ins.MiniProposal {
			ins.MiniProposal = req.GetProposalNum()
			return tx.Model(LogEntry{}).Where("id=?", req.Index).Updates(LogEntry{MiniProposal: req.GetProposalNum()}).Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.PrepareReply{
		AcceptedProposal: ins.AcceptedProposal,
		AcceptedValue:    ins.AcceptedValue,
		NoMoreAccepted:   l.largestAccepted < req.Index,
	}, nil
}

func (l *mysqlLogStore) Accept(req *pb.AcceptRequest) (*pb.AcceptReply, error) {
	ins := LogEntry{}

	err := l.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.FirstOrCreate(&ins, LogEntry{ID: req.Index}).Error; err != nil {
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
		tx.Where(LogEntry{ID: req.Index}).Updates(&ins)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.AcceptReply{
		MiniProposal: ins.MiniProposal,
	}, nil
}
func (l *mysqlLogStore) Learn(req *pb.LearnRequest) (*pb.LearnReply, error) {
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

func (l *mysqlLogStore) PickSlot() int64 {
	return int64(l.picker.Pick())
}
