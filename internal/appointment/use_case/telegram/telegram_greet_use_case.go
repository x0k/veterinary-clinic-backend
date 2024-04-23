package appointment_telegram_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type GreetUseCase[R any] struct {
	greetPresenter appointment.GreetPresenter[R]
}

func NewGreetUseCase[R any](greetPresenter appointment.GreetPresenter[R]) *GreetUseCase[R] {
	return &GreetUseCase[R]{
		greetPresenter: greetPresenter,
	}
}

func (u *GreetUseCase[R]) Greet(ctx context.Context) (R, error) {
	return u.greetPresenter()
}
