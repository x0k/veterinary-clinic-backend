package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrNotFound = errors.New("not found")

type ServicesLoader interface {
	Services(ctx context.Context) ([]shared.Service, error)
}

type ServiceLoader interface {
	Service(ctx context.Context, serviceId shared.ServiceId) (shared.Service, error)
}

type RecordsCreator interface {
	Create(
		ctx context.Context,
		user shared.User,
		service shared.Service,
		dateTime time.Time,
	) (shared.Record, error)
}

type RecordByUserLoader interface {
	RecordByUserId(ctx context.Context, userId shared.UserId) (shared.Record, error)
}

type ActualRecordsLoader interface {
	LoadActualRecords(ctx context.Context, time time.Time) ([]shared.Record, error)
}

type RecordsRemover interface {
	Remove(ctx context.Context, recordId shared.RecordId) error
}

type RecordsArchiver interface {
	ArchiveRecords(ctx context.Context) error
}

type ProductionCalendarLoader interface {
	ProductionCalendar(ctx context.Context) (shared.ProductionCalendar, error)
}

type OpeningHoursLoader interface {
	OpeningHours(ctx context.Context) (shared.OpeningHours, error)
}

type BusyPeriodsLoader interface {
	BusyPeriods(ctx context.Context, t time.Time) (shared.BusyPeriods, error)
}

type WorkBreaksLoader interface {
	WorkBreaks(ctx context.Context) (shared.WorkBreaks, error)
}

type ActualRecordsStateLoader interface {
	ActualRecordsState(ctx context.Context) (shared.ActualRecordsState, error)
}

type ActualRecordsStateSaver interface {
	SaveActualRecordsState(ctx context.Context, state shared.ActualRecordsState) error
}
