package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ClinicRepo interface {
	Services(ctx context.Context) ([]entity.Service, error)
}

type ClinicPresenter[S any] interface {
	RenderServices(services []entity.Service) (S, error)
}

type ClinicUseCase[Services any] struct {
	clinicRepo      ClinicRepo
	clinicPresenter ClinicPresenter[Services]
}

func NewClinicUseCase[Services any](
	clinicRepo ClinicRepo,
	clinicPresenter ClinicPresenter[Services],
) *ClinicUseCase[Services] {
	return &ClinicUseCase[Services]{
		clinicRepo:      clinicRepo,
		clinicPresenter: clinicPresenter,
	}
}

func (u *ClinicUseCase[S]) Services(ctx context.Context) (S, error) {
	services, err := u.clinicRepo.Services(ctx)
	if err != nil {
		return *new(S), err
	}
	return u.clinicPresenter.RenderServices(services)
}
