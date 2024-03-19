package usecase

import "context"

type greetPresenter[R any] interface {
	RenderGreeting() (R, error)
}

type GreetUseCase[R any] struct {
	presenter greetPresenter[R]
}

func NewGreetUseCase[R any](dialogPresenter greetPresenter[R]) *GreetUseCase[R] {
	return &GreetUseCase[R]{
		presenter: dialogPresenter,
	}
}

func (u *GreetUseCase[R]) GreetUser(ctx context.Context) (R, error) {
	return u.presenter.RenderGreeting()
}
