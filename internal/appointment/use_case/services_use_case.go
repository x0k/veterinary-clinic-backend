package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type ServicesUseCase[R any] struct {
	servicesLoader    appointment.ServicesLoader
	servicesPresenter appointment.ServicesPresenter[R]
}

func NewServicesUseCase[R any](
	servicesLoader appointment.ServicesLoader,
	servicesPresenter appointment.ServicesPresenter[R],
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
