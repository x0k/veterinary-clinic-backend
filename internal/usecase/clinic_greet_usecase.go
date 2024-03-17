package usecase

import "context"

type ClinicGreetPresenter[R any] interface {
	RenderGreeting() (R, error)
}

type ClinicGreetUseCase[R any] struct {
	dialogPresenter ClinicGreetPresenter[R]
}

func NewClinicGreetUseCase[R any](dialogPresenter ClinicGreetPresenter[R]) *ClinicGreetUseCase[R] {
	return &ClinicGreetUseCase[R]{
		dialogPresenter: dialogPresenter,
	}
}

func (u *ClinicGreetUseCase[R]) GreetUser(ctx context.Context) (R, error) {
	return u.dialogPresenter.RenderGreeting()
}
