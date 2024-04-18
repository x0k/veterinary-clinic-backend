package appointment_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const scheduleUseCaseName = "appointment_use_case.ScheduleUseCase"

type ScheduleUseCase[R any] struct {
	log               *logger.Logger
	schedulingService *appointment.SchedulingService
	schedulePresenter appointment.SchedulePresenter[R]
	errorPresenter    appointment.ErrorPresenter[R]
}

func NewScheduleUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	schedulePresenter appointment.SchedulePresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *ScheduleUseCase[R] {
	return &ScheduleUseCase[R]{
		log:               log.With(sl.Component(scheduleUseCaseName)),
		schedulingService: schedulingService,
		schedulePresenter: schedulePresenter,
		errorPresenter:    errorPresenter,
	}
}

func (u *ScheduleUseCase[R]) Schedule(ctx context.Context, now, preferredDate time.Time) (R, error) {
	schedule, err := u.schedulingService.Schedule(ctx, now, preferredDate)
	if err != nil {
		u.log.Error(ctx, "failed to get a schedule", sl.Err(err))
		return u.errorPresenter.RenderError(err)
	}
	return u.schedulePresenter.RenderSchedule(now, schedule)
}
