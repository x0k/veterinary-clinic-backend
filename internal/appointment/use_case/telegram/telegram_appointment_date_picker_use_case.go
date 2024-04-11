package appointment_telegram_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type AppointmentDatePickerUseCase[R any] struct {
	schedulingService   appointment.SchedulingService
	datePickerPresenter appointment.DatePickerPresenter[R]
	errorPresenter      appointment.ErrorPresenter[R]
}

func NewAppointmentDatePickerUseCase[R any]() *AppointmentDatePickerUseCase[R] {
	return &AppointmentDatePickerUseCase[R]{}
}

func (u *AppointmentDatePickerUseCase[R]) DatePicker(
	ctx context.Context,
	serviceId appointment.ServiceId,
	now time.Time,
	preferredDate time.Time,
) (R, error) {
	schedule, err := u.schedulingService.Schedule(ctx, now, preferredDate)
	if err != nil {
		return u.errorPresenter.RenderError(err)
	}
	return u.datePickerPresenter.RenderDatePicker(serviceId, schedule)
}
