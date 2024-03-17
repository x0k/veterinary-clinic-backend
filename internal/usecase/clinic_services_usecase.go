package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ClinicServicesRepo interface {
	Services(ctx context.Context) ([]entity.Service, error)
}

type ClinicServicesPresenter[R any] interface {
	RenderServices(services []entity.Service) (R, error)
}

type ClinicServicesUseCase[R any] struct {
	servicesRepo      ClinicServicesRepo
	servicesPresenter ClinicServicesPresenter[R]
}

func NewClinicServices[R any](
	servicesRepo ClinicServicesRepo,
	servicesPresenter ClinicServicesPresenter[R],
) *ClinicServicesUseCase[R] {
	return &ClinicServicesUseCase[R]{
		servicesRepo:      servicesRepo,
		servicesPresenter: servicesPresenter,
	}
}

func (u *ClinicServicesUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.servicesRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.servicesPresenter.RenderServices(services)
}
