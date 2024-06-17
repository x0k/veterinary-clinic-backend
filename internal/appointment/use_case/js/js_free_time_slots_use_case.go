package appointment_js_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const freeTimeSlotsUseCaseName = "appointment_js_use_case.FreeTimeSlotsUseCase"

type FreeTimeSlotsUseCase[R any] struct {
	log                *logger.Logger
	schedulingService  *appointment.SchedulingService
	serviceLoader      appointment.ServiceLoader
	timeSlotsPresenter appointment.TimeSlotsPresenter[R]
	errorPresenter     appointment.ErrorPresenter[R]
}

func NewFreeTimeSlotsUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	serviceLoader appointment.ServiceLoader,
	timeSlotsPresenter appointment.TimeSlotsPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *FreeTimeSlotsUseCase[R] {
	return &FreeTimeSlotsUseCase[R]{
		log:                log.With(sl.Component(freeTimeSlotsUseCaseName)),
		schedulingService:  schedulingService,
		serviceLoader:      serviceLoader,
		timeSlotsPresenter: timeSlotsPresenter,
		errorPresenter:     errorPresenter,
	}
}

func (u *FreeTimeSlotsUseCase[R]) FreeTimeSlots(
	ctx context.Context,
	serviceId appointment.ServiceId,
	now, appointmentDate shared.UTCTime,
) (R, error) {
	service, err := u.serviceLoader(ctx, serviceId)
	if err != nil {
		u.log.Debug(ctx, "failed to load service", sl.Err(err))
		return u.errorPresenter(err)
	}
	sampledFreeTimeSlots, err := u.schedulingService.SampledFreeTimeSlots(
		ctx,
		now,
		appointmentDate,
		service.DurationInMinutes,
	)
	if err != nil {
		u.log.Debug(ctx, "failed to get sampled free time slots", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.timeSlotsPresenter(sampledFreeTimeSlots)
}
