package usecase

import "context"

type clinicGreetPresenter[R any] interface {
	RenderGreeting() (R, error)
}

type ClinicGreetUseCase[R any] struct {
	presenter clinicGreetPresenter[R]
}

func NewClinicGreetUseCase[R any](dialogPresenter clinicGreetPresenter[R]) *ClinicGreetUseCase[R] {
	return &ClinicGreetUseCase[R]{
		presenter: dialogPresenter,
	}
}

func (u *ClinicGreetUseCase[R]) GreetUser(ctx context.Context) (R, error) {
	return u.presenter.RenderGreeting()
}
