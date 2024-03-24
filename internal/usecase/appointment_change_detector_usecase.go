package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrUnexpectedChangeType = errors.New("unexpected change type")
var ErrInvalidRecordUserId = errors.New("invalid record user id")

type appointmentChangePresenter[R any] interface {
	RenderChange(change entity.RecordChange) (R, error)
}

type appointmentChangeActualRecordsStateRepo interface {
	ActualRecordsStateLoader
	ActualRecordsStateSaver
}

type AppointmentChangeDetectorUseCase[R any] struct {
	adminTelegramUserId   entity.UserId
	actualRecordsState    appointmentChangeActualRecordsStateRepo
	recordsRepo           ActualRecordsLoader
	telegramNotifications chan<- entity.NotificationMessage[R]
	presenter             appointmentChangePresenter[R]
}

func NewAppointmentChangeDetectorUseCase[R any](
	adminTelegramUserId entity.UserId,
	actualRecordsState appointmentChangeActualRecordsStateRepo,
	recordsRepo ActualRecordsLoader,
	telegramNotifications chan<- entity.NotificationMessage[R],
	presenter appointmentChangePresenter[R],
) *AppointmentChangeDetectorUseCase[R] {
	return &AppointmentChangeDetectorUseCase[R]{
		adminTelegramUserId:   adminTelegramUserId,
		actualRecordsState:    actualRecordsState,
		recordsRepo:           recordsRepo,
		telegramNotifications: telegramNotifications,
		presenter:             presenter,
	}
}

func (u *AppointmentChangeDetectorUseCase[R]) DetectChanges(
	ctx context.Context,
	now time.Time,
) error {
	state, err := u.actualRecordsState.ActualRecordsState(ctx)
	if err != nil {
		return err
	}
	actualRecords, err := u.recordsRepo.LoadActualRecords(ctx, now)
	if err != nil {
		return err
	}
	changes := state.Update(actualRecords)
	if err := u.actualRecordsState.SaveActualRecordsState(ctx, state); err != nil {
		return err
	}

	errs := make([]error, 0, len(changes))
	for _, change := range changes {
		notification, err := u.presenter.RenderChange(change)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		switch change.Type {
		case entity.RecordCreated:
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  u.adminTelegramUserId,
				Message: notification,
			}
		case entity.RecordStatusChanged:
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  change.Record.User.Id,
				Message: notification,
			}
		case entity.RecordDateTimeChanged:
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  change.Record.User.Id,
				Message: notification,
			}
		case entity.RecordRemoved:
			if change.Record.Status != entity.RecordAwaits {
				continue
			}
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  u.adminTelegramUserId,
				Message: notification,
			}
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  change.Record.User.Id,
				Message: notification,
			}
		default:
			errs = append(errs, fmt.Errorf("%w: %v", ErrUnexpectedChangeType, change.Type))
		}
	}
	return errors.Join(errs...)
}
