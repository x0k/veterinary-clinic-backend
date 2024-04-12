package appointment_telegram_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type AppointmentTimePickerUseCase[R any] struct {
	schedulingService   *appointment.SchedulingService
	serviceLoader       appointment.ServiceLoader
	timePickerPresenter appointment.TimePickerPresenter[R]
	errorPresenter      appointment.ErrorPresenter[R]
}

func NewAppointmentTimePickerUseCase[R any](
	schedulingService *appointment.SchedulingService,
	serviceLoader appointment.ServiceLoader,
	timePickerPresenter appointment.TimePickerPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *AppointmentTimePickerUseCase[R] {
	return &AppointmentTimePickerUseCase[R]{
		schedulingService:   schedulingService,
		serviceLoader:       serviceLoader,
		timePickerPresenter: timePickerPresenter,
		errorPresenter:      errorPresenter,
	}
}

func (u *AppointmentTimePickerUseCase[R]) TimePicker(
	ctx context.Context,
	serviceId appointment.ServiceId,
	now time.Time,
	appointmentDate time.Time,
) (R, error) {
	service, err := u.serviceLoader.Service(ctx, serviceId)
	if err != nil {
		return u.errorPresenter.RenderError(err)
	}
	sampledFreeTimeSlots, err := u.schedulingService.SampledFreeTimeSlots(
		ctx,
		now,
		appointmentDate,
		service.DurationInMinutes,
	)
	if err != nil {
		return u.errorPresenter.RenderError(err)
	}
	return u.timePickerPresenter.RenderTimePicker(serviceId, appointmentDate, sampledFreeTimeSlots)
}
