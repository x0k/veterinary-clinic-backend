package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ServicesPresenter[R any] interface {
	RenderServices(services []ServiceEntity) (R, error)
}

type SchedulePresenter[R any] interface {
	RenderSchedule(now time.Time, schedule Schedule) (R, error)
}

type ErrorPresenter[R any] interface {
	RenderError(err error) (R, error)
}

type RegistrationPresenter[R any] interface {
	RenderRegistration(telegramUserId shared.TelegramUserId) (R, error)
}

type SuccessRegistrationPresenter[R any] interface {
	RenderSuccessRegistration(services []ServiceEntity) (R, error)
}

type ServicesPickerPresenter[R any] interface {
	RenderServicesList(services []ServiceEntity) (R, error)
}

type DatePickerPresenter[R any] interface {
	RenderDatePicker(now time.Time, serviceId ServiceId, schedule Schedule) (R, error)
}

type GreetPresenter[R any] interface {
	RenderGreeting() (R, error)
}

type TimePickerPresenter[R any] interface {
	RenderTimePicker(serviceId ServiceId, appointmentDate time.Time, slots SampledFreeTimeSlots) (R, error)
}

type AppointmentConfirmationPresenter[R any] interface {
	RenderConfirmation(service ServiceEntity, appointmentDateTime time.Time) (R, error)
}

type AppointmentInfoPresenter[R any] interface {
	RenderInfo(
		appointment RecordEntity,
		service ServiceEntity,
	) (R, error)
}

type AppointmentCancelPresenter[R any] interface {
	RenderCancel() (R, error)
}

type EventPresenter[E Event, R any] func(E) (R, error)
