package appointment

import "time"

type ServicesPresenter[R any] interface {
	RenderServices(services []ServiceEntity) (R, error)
}

type SchedulePresenter[R any] interface {
	RenderSchedule(now time.Time, schedule Schedule) (R, error)
}

type ErrorPresenter[R any] interface {
	RenderError(err error) (R, error)
}
