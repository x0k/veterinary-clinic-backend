package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

var ErrUnknownDialog = errors.New("unknown dialog")
var ErrScheduleCalculationFailed = errors.New("schedule calculation failed")
var ErrScheduleRenderingFailed = errors.New("schedule rendering failed")

type DialogRepo interface{}

type DialogPresenter[R any] interface {
	RenderGreeting() (R, error)
	RenderDatePicker() (R, error)
	RenderSchedule([]entity.TimePeriod) (R, error)
	RenderError(error) (R, error)
}

type ClinicDialogUseCase[R any] struct {
	log             *logger.Logger
	dialogPresenter DialogPresenter[R]
	messages        chan entity.DialogMessage[R]
}

func (u *ClinicDialogUseCase[R]) send(dialogId entity.DialogId, message R) {
	u.messages <- entity.DialogMessage[R]{DialogId: dialogId, Message: message}
}

func (u *ClinicDialogUseCase[R]) sendError(ctx context.Context, dialogId entity.DialogId, err error) {
	msg, err := u.dialogPresenter.RenderError(err)
	if err != nil {
		u.log.Error(ctx, "failed to render error", sl.Err(err))
		return
	}
	u.send(dialogId, msg)
}

func NewClinicDialogUseCase[R any](
	log *logger.Logger,
	dialogPresenter DialogPresenter[R],
) *ClinicDialogUseCase[R] {
	return &ClinicDialogUseCase[R]{
		log:             log.With(slog.String("component", "usecase.clinic_dialog.ClinicDialogUseCase")),
		dialogPresenter: dialogPresenter,
		messages:        make(chan entity.DialogMessage[R]),
	}
}

func (u *ClinicDialogUseCase[R]) Messages() <-chan entity.DialogMessage[R] {
	return u.messages
}

func (u *ClinicDialogUseCase[R]) GreetUser(ctx context.Context) (R, error) {
	return u.dialogPresenter.RenderGreeting()
}

func (u *ClinicDialogUseCase[R]) StartScheduleDialog(ctx context.Context) (R, error) {
	return u.dialogPresenter.RenderDatePicker()
}

func (u *ClinicDialogUseCase[R]) FinishScheduleDialog(
	ctx context.Context,
	dialog entity.Dialog,
	t time.Time,
) {
	schedule, err := u.scheduleCalculator.Calculate(t)
	if err != nil {
		u.log.Error(ctx, "failed to calculate schedule", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrScheduleCalculationFailed)
		return
	}
	msg, err := u.dialogPresenter.RenderSchedule(schedule)
	if err != nil {
		u.log.Error(ctx, "failed to render schedule", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrScheduleRenderingFailed)
		return
	}
	u.send(dialog.Id, msg)
}
