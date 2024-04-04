package appointment

import "context"

type ServicesUseCase[R any] struct {
	servicesLoader    ServicesLoader
	servicesPresenter ServicesPresenter[R]
}

func NewServicesUseCase[R any](
	servicesLoader ServicesLoader,
	servicesPresenter ServicesPresenter[R],
) *ServicesUseCase[R] {
	return &ServicesUseCase[R]{
		servicesLoader:    servicesLoader,
		servicesPresenter: servicesPresenter,
	}
}

func (u *ServicesUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.servicesLoader.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.servicesPresenter.RenderServices(services)
}
