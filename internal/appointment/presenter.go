package appointment

type ServicesPresenter[R any] interface {
	RenderServices(services []ServiceEntity) (R, error)
}
