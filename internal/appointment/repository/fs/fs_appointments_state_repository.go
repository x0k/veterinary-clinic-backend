package appointment_fs_repository

import (
	"context"
	"encoding/gob"
	"io"
	"os"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type AppointmentsStateRepository struct {
	name                  string
	mu                    sync.Mutex
	filePath              string
	file                  *os.File
	lastAppointmentsCount int
}

func NewAppointmentsStateRepository(
	name string,
	filePath string,
) *AppointmentsStateRepository {
	return &AppointmentsStateRepository{
		name:     name,
		filePath: filePath,
	}
}

func (r *AppointmentsStateRepository) Name() string {
	return r.name
}

func (r *AppointmentsStateRepository) Start(ctx context.Context) (err error) {
	r.file, err = os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	<-ctx.Done()
	if err := r.file.Sync(); err != nil {
		return err
	}
	return r.file.Close()
}

func (r *AppointmentsStateRepository) AppointmentsState(
	ctx context.Context,
) (appointment.AppointmentsState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	records := make(map[appointment.RecordId]appointment.RecordEntity, r.lastAppointmentsCount)
	if err := gob.NewDecoder(r.file).Decode(&records); err != nil && err != io.EOF {
		return appointment.AppointmentsState{}, err
	}
	if _, err := r.file.Seek(0, 0); err != nil {
		return appointment.AppointmentsState{}, err
	}
	return appointment.NewAppointmentsState(records), nil
}

func (r *AppointmentsStateRepository) SaveAppointmentsState(
	ctx context.Context,
	appointmentsState appointment.AppointmentsState,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.file.Truncate(0); err != nil {
		return err
	}
	encoder := gob.NewEncoder(r.file)
	records := appointmentsState.Appointments()
	r.lastAppointmentsCount = len(records)
	if err := encoder.Encode(records); err != nil {
		return err
	}
	if _, err := r.file.Seek(0, 0); err != nil {
		return err
	}
	return r.file.Sync()
}
