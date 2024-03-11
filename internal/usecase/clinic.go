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

type ClinicUseCase[S any] struct {
	repo ClinicRepo
	p    ClinicPresenter[S]
}

func NewClinicUseCase[S any](repo ClinicRepo, p ClinicPresenter[S]) *ClinicUseCase[S] {
	return &ClinicUseCase[S]{
		repo: repo,
		p:    p,
	}
}

func (u *ClinicUseCase[S]) Services(ctx context.Context) (S, error) {
	services, err := u.repo.Services(ctx)
	if err != nil {
		return *new(S), err
	}
	return u.p.RenderServices(services)
}
