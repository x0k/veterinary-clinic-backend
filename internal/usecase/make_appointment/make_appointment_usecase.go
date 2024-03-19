package make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type makeAppointmentPresenter[R any] interface {
	RenderInfo(appointment entity.Record, service entity.Service) (R, error)
}

type MakeAppointmentUseCase[R any] struct {
	recordsRepo  usecase.RecordsCreator
	servicesRepo usecase.ServiceLoader
	presenter    makeAppointmentPresenter[R]
}

func NewMakeAppointmentUseCase[R any](
	recordsRepo usecase.RecordsCreator,
	servicesRepo usecase.ServiceLoader,
	presenter makeAppointmentPresenter[R],
) *MakeAppointmentUseCase[R] {
	return &MakeAppointmentUseCase[R]{
		recordsRepo:  recordsRepo,
		servicesRepo: servicesRepo,
		presenter:    presenter,
	}
}

func (u *MakeAppointmentUseCase[R]) Make(
	ctx context.Context,
	user entity.User,
	serviceId entity.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	service, err := u.servicesRepo.Service(ctx, serviceId)
	if err != nil {
		return *new(R), err
	}
	record, err := u.recordsRepo.Create(
		ctx,
		user,
		service,
		appointmentDateTime,
	)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderInfo(record, service)
}
