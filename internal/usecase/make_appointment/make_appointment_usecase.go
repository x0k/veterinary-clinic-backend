package make_appointment

import (
	"context"
	"errors"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type appointmentInfoPresenter[R any] interface {
	RenderInfo(appointment entity.Record) (R, error)
}

type makeAppointmentRecordsRepo interface {
	usecase.RecordsCreator
	usecase.RecordByUserLoader
}

type MakeAppointmentUseCase[R any] struct {
	recordsRepo  makeAppointmentRecordsRepo
	servicesRepo usecase.ServiceLoader
	presenter    appointmentInfoPresenter[R]
}

func NewMakeAppointmentUseCase[R any](
	recordsRepo makeAppointmentRecordsRepo,
	servicesRepo usecase.ServiceLoader,
	presenter appointmentInfoPresenter[R],
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
	if rec, err := u.recordsRepo.RecordByUserId(ctx, user.Id); !errors.Is(err, usecase.ErrNotFound) {
		if err != nil {
			return *new(R), err
		}
		return u.presenter.RenderInfo(rec)
	}
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
	return u.presenter.RenderInfo(record)
}
