package appointment_use_case

import (
	"context"

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
		log:               log.With(sl.Component(servicesUseCaseName)),
		servicesLoader:    servicesLoader,
		servicesPresenter: servicesPresenter,
		errorPresenter:    errorPresenter,
	}
}

func (u *ServicesUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.servicesLoader(ctx)
	if err != nil {
		u.log.Debug(ctx, "failed to load services", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.servicesPresenter(services)
}
