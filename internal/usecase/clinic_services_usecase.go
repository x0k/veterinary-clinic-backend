package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type clinicServicesPresenter[R any] interface {
	RenderServices(services []entity.Service) (R, error)
}

type ClinicServicesUseCase[R any] struct {
	servicesRepo ClinicServicesLoader
	presenter    clinicServicesPresenter[R]
}

func NewClinicServicesUseCase[R any](
	servicesRepo ClinicServicesLoader,
	servicesPresenter clinicServicesPresenter[R],
) *ClinicServicesUseCase[R] {
	return &ClinicServicesUseCase[R]{
		servicesRepo: servicesRepo,
		presenter:    servicesPresenter,
	}
}

func (u *ClinicServicesUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.servicesRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderServices(services)
}
