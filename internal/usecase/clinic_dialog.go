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
var ErrLoadingOpeningHoursFailed = errors.New("loading opening hours failed")
var ErrLoadingProductionCalendarFailed = errors.New("loading production calendar failed")
var ErrLoadingWorkBreaksFailed = errors.New("loading work breaks failed")
var ErrLoadingBusyPeriodsFailed = errors.New("loading busy periods failed")

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
	RenderSchedule([]entity.TimePeriod) (R, error)
	RenderError(error) (R, error)
}

type OpeningHoursRepo interface {
	GetOpeningHours(ctx context.Context) (entity.OpeningHours, error)
}

type ProductionCalendarRepo interface {
	GetProductionCalendar(ctx context.Context) (entity.ProductionCalendar, error)
}

type WorkBreaksRepo interface {
	GetWorkBreaks(ctx context.Context, t time.Time) (entity.WorkBreaks, error)
}

type BusyPeriodsRepo interface {
	GetBusyPeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error)
}

type ClinicDialogUseCase[R any] struct {
	log             *logger.Logger
	dialogPresenter DialogPresenter[R]
	messages        chan entity.DialogMessage[R]

	openingHoursRepo       OpeningHoursRepo
	productionCalendarRepo ProductionCalendarRepo
	workBreaksRepo         WorkBreaksRepo
	busyPeriodsRepo        BusyPeriodsRepo
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
	openingHoursRepo OpeningHoursRepo,
	productionCalendarRepo ProductionCalendarRepo,
	workBreaksRepo WorkBreaksRepo,
	busyPeriodsRepo BusyPeriodsRepo,
) *ClinicDialogUseCase[R] {
	return &ClinicDialogUseCase[R]{
		log:                    log.With(slog.String("component", "usecase.clinic_dialog.ClinicDialogUseCase")),
		dialogPresenter:        dialogPresenter,
		messages:               make(chan entity.DialogMessage[R]),
		openingHoursRepo:       openingHoursRepo,
		productionCalendarRepo: productionCalendarRepo,
		workBreaksRepo:         workBreaksRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
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
	openingHours, err := u.openingHoursRepo.GetOpeningHours(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to get opening hours", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingOpeningHoursFailed)
		return
	}
	productionCalendar, err := u.productionCalendarRepo.GetProductionCalendar(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to get production calendar", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingProductionCalendarFailed)
		return
	}

	workBreaks, err := u.workBreaksRepo.GetWorkBreaks(ctx, t)
	if err != nil {
		u.log.Error(ctx, "failed to get work breaks", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingWorkBreaksFailed)
		return
	}
	busyPeriods, err := u.busyPeriodsRepo.GetBusyPeriods(ctx, t)
	if err != nil {
		u.log.Error(ctx, "failed to get busy periods", sl.Err(err))
		u.sendError(ctx, dialog.Id, ErrLoadingBusyPeriodsFailed)
		return
	}
	allBusyPeriods := make([]TitledTimePeriod, 0, len(busyPeriods)+len(workBreaks))
	for _, workBreak := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, TitledTimePeriod{
			TimePeriod: workBreak.Period,
			Title:      workBreak.Title,
			Type:       BusyPeriod,
		})
	}
	for _, busyPeriod := range busyPeriods {
		allBusyPeriods = append(allBusyPeriods, TitledTimePeriod{
			TimePeriod: busyPeriod,
			Title:      "Занято",
			Type:       BusyPeriod,
		})
	}

	now := time.Now()
	dateTime := entity.GoTimeToDateTime(now)
	freePeriodsCalculator := entity.NewFreePeriodsCalculator(
		openingHours,
		productionCalendar,
		dateTime,
	)
	workBreaksCalculator := entity.NewWorkBreaksCalculator(allWorkBreaks)

	entity.NewBusyPeriodsCalculator()

	schedule, err := freePeriodsCalculator.Calculate(t)
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
