package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type clinicServicesLoader interface {
	Services(ctx context.Context) ([]entity.Service, error)
}

type clinicServiceLoader interface {
	Load(ctx context.Context, serviceId entity.ServiceId) (entity.Service, error)
}

type clinicRecordsCreator interface {
	Create(
		ctx context.Context,
		user entity.User,
		serviceId entity.ServiceId,
		dateTime time.Time,
	) (entity.Record, error)
}

type clinicRecordsChecker interface {
	Exists(ctx context.Context, userId entity.UserId) (bool, error)
}

type productionCalendarRepo interface {
	ProductionCalendar(ctx context.Context) (entity.ProductionCalendar, error)
}

type openingHoursRepo interface {
	OpeningHours(ctx context.Context) (entity.OpeningHours, error)
}

type busyPeriodsRepo interface {
	BusyPeriods(ctx context.Context, t time.Time) (entity.BusyPeriods, error)
}

type workBreaksRepo interface {
	WorkBreaks(ctx context.Context) (entity.WorkBreaks, error)
}
