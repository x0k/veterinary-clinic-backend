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

type RecordsRemover interface {
	Remove(ctx context.Context, recordId entity.RecordId) error
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
