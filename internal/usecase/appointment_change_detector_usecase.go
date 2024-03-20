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

type AppointmentChangeDetectorUseCase[R any] struct {
	state map[entity.RecordId]entity.Record

	adminTelegramUserId   entity.UserId
	recordsRepo           ActualRecordsLoader
	telegramNotifications chan<- entity.NotificationMessage[R]
	presenter             appointmentChangePresenter[R]
}

func NewAppointmentChangeDetectorUseCase[R any](
	adminTelegramUserId entity.UserId,
	recordsRepo ActualRecordsLoader,
	telegramNotifications chan<- entity.NotificationMessage[R],
	presenter appointmentChangePresenter[R],
) *AppointmentChangeDetectorUseCase[R] {
	return &AppointmentChangeDetectorUseCase[R]{
		state:                 make(map[entity.RecordId]entity.Record),
		adminTelegramUserId:   adminTelegramUserId,
		recordsRepo:           recordsRepo,
		telegramNotifications: telegramNotifications,
		presenter:             presenter,
	}
}

func (u *AppointmentChangeDetectorUseCase[R]) DetectChanges(
	ctx context.Context,
	now time.Time,
) error {
	actualRecords, err := u.recordsRepo.LoadActualRecords(ctx, now)
	if err != nil {
		return err
	}
	changes := entity.DetectChanges(u.state, actualRecords)
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
			if change.Record.UserId == nil {
				errs = append(errs, ErrInvalidRecordUserId)
				continue
			}
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  *change.Record.UserId,
				Message: notification,
			}
		case entity.RecordDateTimeChanged:
			if change.Record.UserId == nil {
				errs = append(errs, ErrInvalidRecordUserId)
				continue
			}
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  *change.Record.UserId,
				Message: notification,
			}
		case entity.RecordRemoved:
			if change.Record.Status == entity.RecordInWork {
				continue
			}
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  u.adminTelegramUserId,
				Message: notification,
			}
			if change.Record.UserId == nil {
				errs = append(errs, ErrInvalidRecordUserId)
				continue
			}
			u.telegramNotifications <- entity.NotificationMessage[R]{
				UserId:  *change.Record.UserId,
				Message: notification,
			}
		default:
			errs = append(errs, fmt.Errorf("%w: %v", ErrUnexpectedChangeType, change.Type))
		}
	}
	return errors.Join(errs...)
}
