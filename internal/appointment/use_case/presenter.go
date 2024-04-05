package appointment_use_case

import "github.com/x0k/veterinary-clinic-backend/internal/appointment"

type ServicesPresenter[R any] interface {
	RenderServices(services []appointment.ServiceEntity) (R, error)
}

type SchedulePresenter[R any] interface {
	RenderSchedule(schedule appointment.Schedule) (R, error)
}
