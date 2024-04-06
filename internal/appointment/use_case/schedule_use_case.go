package appointment_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type ScheduleUseCase[R any] struct {
	schedulingService appointment.SchedulingService
	schedulePresenter appointment.SchedulePresenter[R]
	errorPresenter    appointment.ErrorPresenter[R]
}

func (u *ScheduleUseCase[R]) Schedule(ctx context.Context, now, preferredDate time.Time) (R, error) {
	schedule, err := u.schedulingService.Schedule(ctx, now, preferredDate)
	if err != nil {
		return u.errorPresenter.RenderError(err)
	}
	return u.schedulePresenter.RenderSchedule(schedule)
}
