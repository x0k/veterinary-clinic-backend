package appointment_telegram_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const appointmentDatePickerUseCaseName = "appointment_telegram_use_case.AppointmentDatePickerUseCase"

type AppointmentDatePickerUseCase[R any] struct {
	log                 *logger.Logger
	schedulingService   *appointment.SchedulingService
	datePickerPresenter appointment.DatePickerPresenter[R]
	errorPresenter      appointment.ErrorPresenter[R]
}

func NewAppointmentDatePickerUseCase[R any](
	log *logger.Logger,
	schedulingService *appointment.SchedulingService,
	datePickerPresenter appointment.DatePickerPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *AppointmentDatePickerUseCase[R] {
	return &AppointmentDatePickerUseCase[R]{
		log:                 log.With(sl.Component(appointmentDatePickerUseCaseName)),
		schedulingService:   schedulingService,
		datePickerPresenter: datePickerPresenter,
		errorPresenter:      errorPresenter,
	}
}

func (u *AppointmentDatePickerUseCase[R]) DatePicker(
	ctx context.Context,
	serviceId appointment.ServiceId,
	now time.Time,
	preferredDate time.Time,
) (R, error) {
	schedule, err := u.schedulingService.Schedule(ctx, now, preferredDate)
	if err != nil {
		u.log.Error(ctx, "failed to get a schedule", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.datePickerPresenter(now, serviceId, schedule)
}
