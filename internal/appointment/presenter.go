package appointment

type ServicesPresenter[R any] interface {
	RenderServices(services []ServiceEntity) (R, error)
}

type SchedulePresenter[R any] interface {
	RenderSchedule(schedule Schedule) (R, error)
}

type ErrorPresenter[R any] interface {
	RenderError(err error) (R, error)
}
