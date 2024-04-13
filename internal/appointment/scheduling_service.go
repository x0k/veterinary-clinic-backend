package appointment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrInvalidRecordId = errors.New("invalid record id")
var ErrPeriodIsLocked = errors.New("periods is locked")
var ErrDateTimePeriodIsOccupied = errors.New("date time period is occupied")
var ErrAnotherAppointmentIsAlreadyScheduled = errors.New("another appointment is already scheduled")
var ErrInvalidAppointmentStatusForCancel = errors.New("invalid appointment status")

type SchedulingService struct {
	log       *logger.Logger
	periodsMu sync.Mutex
	periods   []shared.DateTimePeriod

	sampleRateInMinutes             SampleRateInMinutes
	appointmentCreator              AppointmentCreator
	productionCalendarLoader        ProductionCalendarLoader
	workingHoursLoader              WorkingHoursLoader
	busyPeriodsLoader               BusyPeriodsLoader
	workBreaksLoader                WorkBreaksLoader
	customerActiveAppointmentLoader CustomerActiveAppointmentLoader
	appointmentRemover              AppointmentRemover
}

func NewSchedulingService(
	log *logger.Logger,
	sampleRateInMinutes SampleRateInMinutes,
	appointmentCreator AppointmentCreator,
	productionCalendarLoader ProductionCalendarLoader,
	workingHoursLoader WorkingHoursLoader,
	busyPeriodsLoader BusyPeriodsLoader,
	workBreaksLoader WorkBreaksLoader,
	customerActiveAppointmentLoader CustomerActiveAppointmentLoader,
	appointmentRemover AppointmentRemover,
) *SchedulingService {
	return &SchedulingService{
		log:                             log.With(slog.String("component", "SchedulingService")),
		sampleRateInMinutes:             sampleRateInMinutes,
		appointmentCreator:              appointmentCreator,
		productionCalendarLoader:        productionCalendarLoader,
		workingHoursLoader:              workingHoursLoader,
		busyPeriodsLoader:               busyPeriodsLoader,
		workBreaksLoader:                workBreaksLoader,
		customerActiveAppointmentLoader: customerActiveAppointmentLoader,
		appointmentRemover:              appointmentRemover,
	}
}

func (s *SchedulingService) lockPeriod(period shared.DateTimePeriod) error {
	s.periodsMu.Lock()
	defer s.periodsMu.Unlock()
	for _, p := range s.periods {
		if shared.DateTimePeriodApi.IsValidPeriod(
			shared.DateTimePeriodApi.IntersectPeriods(p, period),
		) {
			return fmt.Errorf("%w: %s", ErrPeriodIsLocked, period)
		}
	}
	s.periods = append(s.periods, period)
	return nil
}

func (s *SchedulingService) unLockPeriod(period shared.DateTimePeriod) error {
	s.periodsMu.Lock()
	defer s.periodsMu.Unlock()
	index := slices.Index(s.periods, period)
	if index == -1 {
		return nil
	}
	s.periods = slices.Delete(s.periods, index, index+1)
	return nil
}

func (s *SchedulingService) MakeAppointment(
	ctx context.Context,
	now time.Time,
	appointmentDate time.Time,
	customer CustomerEntity,
	service ServiceEntity,
) (AppointmentAggregate, error) {
	appointmentDateTime := shared.GoTimeToDateTime(appointmentDate)
	dateTimePeriod := shared.DateTimePeriod{
		Start: appointmentDateTime,
		End: shared.DateTime{
			Date: appointmentDateTime.Date,
			Time: shared.MakeTimeShifter(shared.Time{
				Minutes: service.DurationInMinutes.Minutes(),
			})(appointmentDateTime.Time),
		},
	}
	if err := s.lockPeriod(dateTimePeriod); err != nil {
		return AppointmentAggregate{}, err
	}
	defer func() {
		if err := s.unLockPeriod(dateTimePeriod); err != nil {
			s.log.Error(ctx, "failed to unlock period", sl.Err(err))
		}
	}()
	existedAppointment, err := s.customerActiveAppointmentLoader.CustomerActiveAppointment(ctx, customer)
	if !errors.Is(err, shared.ErrNotFound) {
		if err != nil {
			return AppointmentAggregate{}, err
		}
		return AppointmentAggregate{}, fmt.Errorf("%w: %s", ErrAnotherAppointmentIsAlreadyScheduled, existedAppointment.Id())
	}
	productionCalendar, err := s.productionCalendar(ctx)
	if err != nil {
		return AppointmentAggregate{}, err
	}
	busyPeriods, err := s.busyPeriodsLoader.BusyPeriods(ctx, appointmentDate)
	if err != nil {
		return AppointmentAggregate{}, err
	}
	datWorkBreaks, err := s.dayWorkBreaks(ctx, appointmentDate)
	if err != nil {
		return AppointmentAggregate{}, err
	}
	freeTimeSlots, err := s.freeTimeSlots(
		ctx,
		now,
		appointmentDate,
		productionCalendar,
		busyPeriods,
		datWorkBreaks,
	)
	if err != nil {
		return AppointmentAggregate{}, err
	}
	if !freeTimeSlots.Includes(shared.TimePeriod{
		Start: dateTimePeriod.Start.Time,
		End:   dateTimePeriod.End.Time,
	}) {
		return AppointmentAggregate{}, fmt.Errorf("%w: %s", ErrDateTimePeriodIsOccupied, dateTimePeriod)
	}
	record, err := NewRecord(
		TemporalRecordId,
		RecordAwaits,
		false,
		dateTimePeriod,
		customer.Id,
		service.Id,
		now,
	)
	if err != nil {
		return AppointmentAggregate{}, err
	}
	app, err := NewAppointmentAggregate(record, service, customer)
	if err != nil {
		return AppointmentAggregate{}, err
	}
	if err := s.appointmentCreator.CreateAppointment(ctx, &app); err != nil {
		return AppointmentAggregate{}, err
	}
	if app.Id() == TemporalRecordId {
		return AppointmentAggregate{}, fmt.Errorf("%w: %s", ErrInvalidRecordId, app.Id())
	}
	return app, nil
}

func (s *SchedulingService) Schedule(
	ctx context.Context,
	now time.Time,
	preferredDate time.Time,
) (Schedule, error) {
	productionCalendar, err := s.productionCalendar(ctx)
	if err != nil {
		return Schedule{}, err
	}
	appointmentDate := productionCalendar.DayOrNextWorkingDay(preferredDate)
	busyPeriods, err := s.busyPeriodsLoader.BusyPeriods(ctx, appointmentDate)
	if err != nil {
		return Schedule{}, err
	}
	dayWorkBreaks, err := s.dayWorkBreaks(ctx, appointmentDate)
	if err != nil {
		return Schedule{}, err
	}
	freeTimeSlots, err := s.freeTimeSlots(
		ctx,
		now,
		appointmentDate,
		productionCalendar,
		busyPeriods,
		dayWorkBreaks,
	)
	if err != nil {
		return Schedule{}, err
	}
	return NewSchedule(
		now,
		appointmentDate,
		productionCalendar,
		freeTimeSlots,
		busyPeriods,
		dayWorkBreaks,
	), nil
}

func (s *SchedulingService) SampledFreeTimeSlots(
	ctx context.Context,
	now time.Time,
	appointmentDate time.Time,
	durationInMinutes shared.DurationInMinutes,
) (SampledFreeTimeSlots, error) {
	productionCalendar, err := s.productionCalendar(ctx)
	if err != nil {
		return SampledFreeTimeSlots{}, err
	}
	busyPeriods, err := s.busyPeriodsLoader.BusyPeriods(ctx, appointmentDate)
	if err != nil {
		return SampledFreeTimeSlots{}, err
	}
	dayWorkBreaks, err := s.dayWorkBreaks(ctx, appointmentDate)
	if err != nil {
		return SampledFreeTimeSlots{}, err
	}
	freeTimeSlots, err := s.freeTimeSlots(
		ctx,
		now,
		appointmentDate,
		productionCalendar,
		busyPeriods,
		dayWorkBreaks,
	)
	if err != nil {
		return SampledFreeTimeSlots{}, err
	}
	return NewSampleFreeTimeSlots(
		durationInMinutes,
		s.sampleRateInMinutes,
		freeTimeSlots,
	), nil
}

func (s *SchedulingService) CancelAppointmentForCustomer(
	ctx context.Context,
	customer CustomerEntity,
) error {
	app, err := s.customerActiveAppointmentLoader.CustomerActiveAppointment(ctx, customer)
	if err != nil {
		return err
	}
	if app.Status() != RecordAwaits {
		return fmt.Errorf("%w: %s", ErrInvalidAppointmentStatusForCancel, app.Status())
	}
	return s.appointmentRemover.RemoveAppointment(ctx, app.Id())
}

func (s *SchedulingService) productionCalendar(ctx context.Context) (ProductionCalendar, error) {
	pc, err := s.productionCalendarLoader.ProductionCalendar(ctx)
	if err != nil {
		return nil, err
	}
	return pc.WithoutSaturdayWeekend(), nil
}

func (s *SchedulingService) dayWorkBreaks(ctx context.Context, day time.Time) (DayWorkBreaks, error) {
	workBreaks, err := s.workBreaksLoader.WorkBreaks(ctx)
	if err != nil {
		return DayWorkBreaks{}, err
	}
	return workBreaks.ForDay(day)
}

func (s *SchedulingService) freeTimeSlots(
	ctx context.Context,
	now time.Time,
	appointmentDate time.Time,
	productionCalendar ProductionCalendar,
	busyPeriods BusyPeriods,
	dayWorkBreaks DayWorkBreaks,
) (FreeTimeSlots, error) {
	workingHours, err := s.workingHoursLoader.WorkingHours(ctx)
	if err != nil {
		return FreeTimeSlots{}, err
	}
	datTimePeriods, err := workingHours.ForDay(appointmentDate).
		OmitPast(shared.GoTimeToDateTime(now)).
		ConsiderProductionCalendar(productionCalendar)
	if err != nil {
		return FreeTimeSlots{}, err
	}
	return NewFreeTimeSlots(
		datTimePeriods,
		busyPeriods,
		dayWorkBreaks,
	)
}
