package make_appointment

import (
	"context"
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type servicePickerPresenter[R any] interface {
	RenderServicesList(services []entity.Service) (R, error)
}

type ServicePickerUseCase[R any] struct {
	servicesRepo             usecase.ServicesLoader
	recordsRepo              usecase.RecordByUserLoader
	servicePickerPresenter   servicePickerPresenter[R]
	appointmentInfoPresenter appointmentInfoPresenter[R]
}

func NewServicePickerUseCase[R any](
	servicesRepo usecase.ServicesLoader,
	recordsRepo usecase.RecordByUserLoader,
	servicePickerPresenter servicePickerPresenter[R],
	appointmentInfoPresenter appointmentInfoPresenter[R],
) *ServicePickerUseCase[R] {
	return &ServicePickerUseCase[R]{
		servicesRepo:             servicesRepo,
		recordsRepo:              recordsRepo,
		servicePickerPresenter:   servicePickerPresenter,
		appointmentInfoPresenter: appointmentInfoPresenter,
	}
}

func (u *ServicePickerUseCase[R]) ServicesPicker(ctx context.Context, userId entity.UserId) (R, error) {
	if record, err := u.recordsRepo.RecordByUserId(ctx, userId); !errors.Is(err, usecase.ErrNotFound) {
		if err != nil {
			return *new(R), err
		}
		return u.appointmentInfoPresenter.RenderInfo(record)
	}
	services, err := u.servicesRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.servicePickerPresenter.RenderServicesList(services)
}
