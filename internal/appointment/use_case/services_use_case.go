package appointment_use_case

import (
	"context"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const servicesUseCaseName = "appointment_use_case.ServicesUseCase"

type ServicesUseCase[R any] struct {
	log               *logger.Logger
	servicesLoader    appointment.ServicesLoader
	servicesPresenter appointment.ServicesPresenter[R]
	errorPresenter    appointment.ErrorPresenter[R]
}

func NewServicesUseCase[R any](
	log *logger.Logger,
	servicesLoader appointment.ServicesLoader,
	servicesPresenter appointment.ServicesPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *ServicesUseCase[R] {
	return &ServicesUseCase[R]{
		log:               log.With(slog.String("component", servicesUseCaseName)),
		servicesLoader:    servicesLoader,
		servicesPresenter: servicesPresenter,
		errorPresenter:    errorPresenter,
	}
}

func (u *ServicesUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.servicesLoader.Services(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to load services", sl.Err(err))
		return u.errorPresenter.RenderError(err)
	}
	return u.servicesPresenter.RenderServices(services)
}
