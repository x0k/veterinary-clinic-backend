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
var ErrScheduleRenderingFailed = errors.New("schedule rendering failed")
var ErrLoadingWorkBreaksFailed = errors.New("loading work breaks failed")
var ErrLoadingBusyPeriodsFailed = errors.New("loading busy periods failed")
var ErrLoadingFreePeriodsFailed = errors.New("loading free periods failed")
var ErrNextAvailableDayCalculationFailed = errors.New("next available day calculation failed")

type DialogPresenter[R any] interface {
	RenderGreeting() (R, error)
	RenderDatePicker(time.Time) (R, error)
	RenderSchedule(entity.Schedule) (R, error)
	RenderError(error) (R, error)
}

type WorkBreaksRepo interface {
	WorkBreaks(ctx context.Context, t time.Time) (entity.WorkBreaks, error)
}

type BusyPeriodsRepo interface {
	BusyPeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error)
}

type FreePeriodsRepo interface {
	FreePeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error)
	NextAvailableDay(ctx context.Context, t time.Time) (time.Time, error)
}

type ClinicDialogUseCase[R any] struct {
	log             *logger.Logger
	dialogPresenter DialogPresenter[R]
	messages        chan entity.DialogMessage[R]

	workBreaksRepo  WorkBreaksRepo
	busyPeriodsRepo BusyPeriodsRepo
	freePeriodsRepo FreePeriodsRepo
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
	workBreaksRepo WorkBreaksRepo,
	busyPeriodsRepo BusyPeriodsRepo,
	freePeriodsRepo FreePeriodsRepo,
) *ClinicDialogUseCase[R] {
	return &ClinicDialogUseCase[R]{
		log:             log.With(slog.String("component", "usecase.clinic_dialog.ClinicDialogUseCase")),
		dialogPresenter: dialogPresenter,
		messages:        make(chan entity.DialogMessage[R]),
		workBreaksRepo:  workBreaksRepo,
		busyPeriodsRepo: busyPeriodsRepo,
		freePeriodsRepo: freePeriodsRepo,
	}
}

func (u *ClinicDialogUseCase[R]) Messages() <-chan entity.DialogMessage[R] {
	return u.messages
}

func (u *ClinicDialogUseCase[R]) GreetUser(ctx context.Context) (R, error) {
	return u.dialogPresenter.RenderGreeting()
}

func (u *ClinicDialogUseCase[R]) StartScheduleDialog(ctx context.Context) (R, error) {
	t, err := u.freePeriodsRepo.NextAvailableDay(ctx, time.Now())
	if err != nil {
		u.log.Error(ctx, "failed to get next available day", sl.Err(err))
		return *new(R), ErrNextAvailableDayCalculationFailed
	}
	return u.dialogPresenter.RenderDatePicker(t)
}

func (u *ClinicDialogUseCase[R]) FinishScheduleDialog(
	ctx context.Context,
	dialog entity.Dialog,
	t time.Time,
) {
	freePeriods, err := u.freePeriodsRepo.FreePeriods(ctx, t)
	if err != nil {
		u.log.Error(ctx, "failed to get free periods", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingFreePeriodsFailed)
		return
	}
	workBreaks, err := u.workBreaksRepo.WorkBreaks(ctx, t)
	if err != nil {
		u.log.Error(ctx, "failed to get work breaks", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingWorkBreaksFailed)
		return
	}
	busyPeriods, err := u.busyPeriodsRepo.BusyPeriods(ctx, t)
	if err != nil {
		u.log.Error(ctx, "failed to get busy periods", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingBusyPeriodsFailed)
		return
	}
	schedulePeriods := entity.CalculateSchedulePeriods(freePeriods, busyPeriods, workBreaks)

	msg, err := u.dialogPresenter.RenderSchedule(
		entity.NewSchedule(t, schedulePeriods),
	)
	if err != nil {
		u.log.Error(ctx, "failed to render schedule", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrScheduleRenderingFailed)
		return
	}
	u.send(dialog.Id, msg)
}
