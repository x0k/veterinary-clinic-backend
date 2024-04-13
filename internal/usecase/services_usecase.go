package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type servicesPresenter[R any] interface {
	RenderServices(services []shared.Service) (R, error)
}

type ServicesUseCase[R any] struct {
	servicesRepo ServicesLoader
	presenter    servicesPresenter[R]
}

func NewServicesUseCase[R any](
	servicesRepo ServicesLoader,
	servicesPresenter servicesPresenter[R],
) *ServicesUseCase[R] {
	return &ServicesUseCase[R]{
		servicesRepo: servicesRepo,
		presenter:    servicesPresenter,
	}
}

func (u *ServicesUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.servicesRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderServices(services)
}
