package appointment_js_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const dayOrNextWorkingDayUseCaseName = "appointment_js_use_case.DayOrNextWorkingDayUseCase"

type DayOrNextWorkingDayUseCase[R any] struct {
	log                      *logger.Logger
	productionCalendarLoader appointment.ProductionCalendarLoader
	dayPresenter             appointment.DayPresenter[R]
	errorPresenter           appointment.ErrorPresenter[R]
}

func NewDayOrNextWorkingDayUseCase[R any](
	log *logger.Logger,
	productionCalendarLoader appointment.ProductionCalendarLoader,
	dayPresenter appointment.DayPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *DayOrNextWorkingDayUseCase[R] {
	return &DayOrNextWorkingDayUseCase[R]{
		log:                      log.With(sl.Component(dayOrNextWorkingDayUseCaseName)),
		productionCalendarLoader: productionCalendarLoader,
		dayPresenter:             dayPresenter,
		errorPresenter:           errorPresenter,
	}
}

func (u *DayOrNextWorkingDayUseCase[R]) DayOrNextWorkingDay(ctx context.Context, now time.Time) (R, error) {
	cal, err := u.productionCalendarLoader(ctx)
	if err != nil {
		u.log.Debug(ctx, "failed to load production calendar", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.dayPresenter(cal.DayOrNextWorkingDay(now))
}
