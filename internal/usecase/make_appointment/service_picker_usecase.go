package make_appointment

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type servicePickerPresenter[R any] interface {
	RenderServicesList(services []entity.Service) (R, error)
}

type ServicePickerUseCase[R any] struct {
	servicesRepo usecase.ServicesLoader
	presenter    servicePickerPresenter[R]
}

func NewServicePickerUseCase[R any](
	servicesRepo usecase.ServicesLoader,
	presenter servicePickerPresenter[R],
) *ServicePickerUseCase[R] {
	return &ServicePickerUseCase[R]{
		servicesRepo: servicesRepo,
		presenter:    presenter,
	}
}

func (u *ServicePickerUseCase[R]) ServicesPicker(ctx context.Context) (R, error) {
	services, err := u.servicesRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderServicesList(services)
}
