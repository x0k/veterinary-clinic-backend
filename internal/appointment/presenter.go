package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
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
	RenderRegistration(telegramUserId entity.TelegramUserId) (R, error)
}

type SuccessRegistrationPresenter[R any] interface {
	RenderSuccessRegistration(services []ServiceEntity) (R, error)
}

type ServicesPickerPresenter[R any] interface {
	RenderServicesList(services []ServiceEntity) (R, error)
}

type DatePickerPresenter[R any] interface {
	RenderDatePicker(serviceId ServiceId, schedule Schedule) (R, error)
}
