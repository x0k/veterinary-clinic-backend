package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ClinicRepo interface {
	Services(ctx context.Context) ([]entity.Service, error)
}

type ClinicPresenter[R any] interface {
	RenderServices(services []entity.Service) (R, error)
}

type ClinicUseCase[R any] struct {
	clinicRepo      ClinicRepo
	clinicPresenter ClinicPresenter[R]
}

func NewClinicUseCase[R any](
	clinicRepo ClinicRepo,
	clinicPresenter ClinicPresenter[R],
) *ClinicUseCase[R] {
	return &ClinicUseCase[R]{
		clinicRepo:      clinicRepo,
		clinicPresenter: clinicPresenter,
	}
}

func (u *ClinicUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.clinicRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.clinicPresenter.RenderServices(services)
}
