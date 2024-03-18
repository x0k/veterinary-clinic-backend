package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ClinicServicesLoader interface {
	Services(ctx context.Context) ([]entity.Service, error)
}

type ClinicServiceLoader interface {
	Load(ctx context.Context, serviceId entity.ServiceId) (entity.Service, error)
}

type ClinicRecordsCreator interface {
	Create(
		ctx context.Context,
		user entity.User,
		serviceId entity.ServiceId,
		dateTime time.Time,
	) (entity.Record, error)
}

type ClinicRecordsChecker interface {
	Exists(ctx context.Context, userId entity.UserId) (bool, error)
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
