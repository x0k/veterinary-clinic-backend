package usecase

import (
	"context"
	"errors"
	"log/slog"
	"slices"
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

type TimePeriodType int

const (
	FreePeriod TimePeriodType = iota
	BusyPeriod
)

type TitledTimePeriod struct {
	entity.TimePeriod
	Type  TimePeriodType
	Title string
}

type DialogPresenter[R any] interface {
	RenderGreeting() (R, error)
	RenderDatePicker() (R, error)
	RenderSchedule([]TitledTimePeriod) (R, error)
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
	allBusyPeriods := make([]entity.TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}

	actualFreePeriods := entity.TimePeriodApi.SortAndUnitePeriods(
		entity.TimePeriodApi.SubtractPeriodsFromPeriods(
			freePeriods,
			allBusyPeriods,
		),
	)

	schedule := make([]TitledTimePeriod, 0, len(actualFreePeriods)+len(allBusyPeriods))
	for _, p := range actualFreePeriods {
		schedule = append(schedule, TitledTimePeriod{
			TimePeriod: p,
			Type:       FreePeriod,
			Title:      "Свободно",
		})
	}
	for _, p := range busyPeriods {
		schedule = append(schedule, TitledTimePeriod{
			TimePeriod: p,
			Type:       BusyPeriod,
			Title:      "Занято",
		})
	}
	for _, p := range workBreaks {
		schedule = append(schedule, TitledTimePeriod{
			TimePeriod: p.Period,
			Type:       BusyPeriod,
			Title:      p.Title,
		})
	}
	slices.SortFunc(schedule, func(a, b TitledTimePeriod) int {
		return entity.TimePeriodApi.ComparePeriods(a.TimePeriod, b.TimePeriod)
	})

	msg, err := u.dialogPresenter.RenderSchedule(schedule)
	if err != nil {
		u.log.Error(ctx, "failed to render schedule", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrScheduleRenderingFailed)
		return
	}
	u.send(dialog.Id, msg)
}
