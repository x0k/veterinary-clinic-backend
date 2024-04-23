package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ServicesPresenter[R any] func(services []ServiceEntity) (R, error)

type SchedulePresenter[R any] func(now time.Time, schedule Schedule) (R, error)

type ErrorPresenter[R any] func(err error) (R, error)

type RegistrationPresenter[R any] func(telegramUserId shared.TelegramUserId) (R, error)

type SuccessRegistrationPresenter[R any] func(services []ServiceEntity) (R, error)

type ServicesPickerPresenter[R any] func(services []ServiceEntity) (R, error)

type DatePickerPresenter[R any] func(now time.Time, serviceId ServiceId, schedule Schedule) (R, error)

type GreetPresenter[R any] func() (R, error)

type TimePickerPresenter[R any] func(
	serviceId ServiceId,
	appointmentDate time.Time,
	slots SampledFreeTimeSlots,
) (R, error)

type AppointmentConfirmationPresenter[R any] func(
	service ServiceEntity,
	appointmentDateTime time.Time,
) (R, error)

type AppointmentInfoPresenter[R any] func(
	appointment RecordEntity,
	service ServiceEntity,
) (R, error)

type AppointmentCancelPresenter[R any] func() (R, error)

type EventPresenter[E Event, R any] func(E) (R, error)

type ChangedEventPresenter[R any] func(ChangedEvent, CustomerEntity, ServiceEntity) (R, error)
