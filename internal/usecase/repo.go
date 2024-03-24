package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrNotFound = errors.New("not found")

type ServicesLoader interface {
	Services(ctx context.Context) ([]entity.Service, error)
}

type ServiceLoader interface {
	Service(ctx context.Context, serviceId entity.ServiceId) (entity.Service, error)
}

type RecordsCreator interface {
	Create(
		ctx context.Context,
		user entity.User,
		service entity.Service,
		dateTime time.Time,
	) (entity.Record, error)
}

type RecordByUserLoader interface {
	RecordByUserId(ctx context.Context, userId entity.UserId) (entity.Record, error)
}

type ActualRecordsLoader interface {
	LoadActualRecords(ctx context.Context, time time.Time) ([]entity.Record, error)
}

type RecordsRemover interface {
	Remove(ctx context.Context, recordId entity.RecordId) error
}

type RecordsArchiver interface {
	ArchiveRecords(ctx context.Context) error
}

type ProductionCalendarLoader interface {
	ProductionCalendar(ctx context.Context) (entity.ProductionCalendar, error)
}

type OpeningHoursLoader interface {
	OpeningHours(ctx context.Context) (entity.OpeningHours, error)
}

type BusyPeriodsLoader interface {
	BusyPeriods(ctx context.Context, t time.Time) (entity.BusyPeriods, error)
}

type WorkBreaksLoader interface {
	WorkBreaks(ctx context.Context) (entity.WorkBreaks, error)
}

type ActualRecordsStateLoader interface {
	ActualRecordsState(ctx context.Context) (entity.ActualRecordsState, error)
}

type ActualRecordsStateSaver interface {
	SaveActualRecordsState(ctx context.Context, state entity.ActualRecordsState) error
}
