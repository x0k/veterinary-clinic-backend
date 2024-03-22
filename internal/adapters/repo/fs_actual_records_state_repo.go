package repo

import (
	"context"
	"encoding/gob"
	"io"
	"os"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type FsActualRecordsStateRepo struct {
	mu               sync.Mutex
	statePath        string
	file             *os.File
	lastRecordsCount int
}

func NewFsActualRecordsStateRepo(statePath string) *FsActualRecordsStateRepo {
	return &FsActualRecordsStateRepo{statePath: statePath}
}

func (r *FsActualRecordsStateRepo) Start(ctx context.Context) (err error) {
	r.file, err = os.OpenFile(r.statePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	<-ctx.Done()
	if err := r.file.Sync(); err != nil {
		return err
	}
	return r.file.Close()
}

func (s *FsActualRecordsStateRepo) ActualRecordsState(ctx context.Context) (entity.ActualRecordsState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	records := make(map[entity.RecordId]entity.Record, s.lastRecordsCount)
	if err := gob.NewDecoder(s.file).Decode(&records); err != nil && err != io.EOF {
		return entity.ActualRecordsState{}, err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return entity.ActualRecordsState{}, err
	}
	return entity.NewActualRecordsState(records), nil
}

func (s *FsActualRecordsStateRepo) SaveActualRecordsState(ctx context.Context, state entity.ActualRecordsState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.file.Truncate(0); err != nil {
		return err
	}
	encoder := gob.NewEncoder(s.file)
	records := state.Records()
	s.lastRecordsCount = len(records)
	if err := encoder.Encode(records); err != nil {
		return err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	return s.file.Sync()
}
