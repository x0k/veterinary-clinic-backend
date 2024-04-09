package appointment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

var ErrInvalidRecordId = errors.New("invalid record id")
var ErrPeriodIsLocked = errors.New("periods is locked")
var ErrDateTimePeriodIsOccupied = errors.New("date time period is occupied")

type SchedulingService struct {
	log                       *logger.Logger
	periodsMu                 sync.Mutex
	periods                   []entity.DateTimePeriod
	appointmentPeriodsChecker AppointmentPeriodChecker
	appointmentCreator        AppointmentCreator
	productionCalendarLoader  ProductionCalendarLoader
	workingHoursLoader        WorkingHoursLoader
	busyPeriodsLoader         BusyPeriodsLoader
	workBreaksLoader          WorkBreaksLoader
}

func NewSchedulingService(
	log *logger.Logger,
	appointmentPeriodChecker AppointmentPeriodChecker,
	appointmentCreator AppointmentCreator,
	productionCalendarLoader ProductionCalendarLoader,
	workingHoursLoader WorkingHoursLoader,
	busyPeriodsLoader BusyPeriodsLoader,
	workBreaksLoader WorkBreaksLoader,
) *SchedulingService {
	return &SchedulingService{
		log:                       log.With(slog.String("component", "SchedulingService")),
		appointmentPeriodsChecker: appointmentPeriodChecker,
		appointmentCreator:        appointmentCreator,
		productionCalendarLoader:  productionCalendarLoader,
		workingHoursLoader:        workingHoursLoader,
		busyPeriodsLoader:         busyPeriodsLoader,
		workBreaksLoader:          workBreaksLoader,
	}
}

func (s *SchedulingService) lockPeriod(period entity.DateTimePeriod) error {
	s.periodsMu.Lock()
	defer s.periodsMu.Unlock()
	for _, p := range s.periods {
		if entity.DateTimePeriodApi.IsValidPeriod(
			entity.DateTimePeriodApi.IntersectPeriods(p, period),
		) {
			return fmt.Errorf("%w: %s", ErrPeriodIsLocked, period)
		}
	}
	s.periods = append(s.periods, period)
	return nil
}

func (s *SchedulingService) unLockPeriod(period entity.DateTimePeriod) error {
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
	customer CustomerEntity,
	service ServiceEntity,
	dateTimePeriod entity.DateTimePeriod,
) (*AppointmentAggregate, error) {
	if err := s.lockPeriod(dateTimePeriod); err != nil {
		return nil, err
	}
	defer func() {
		if err := s.unLockPeriod(dateTimePeriod); err != nil {
			s.log.Error(ctx, "failed to unlock period", sl.Err(err))
		}
	}()
	// TODO: Check weekends, holidays, etc
	isBusy, err := s.appointmentPeriodsChecker.IsAppointmentPeriodBusy(ctx, dateTimePeriod)
	if err != nil {
		return nil, err
	}
	if isBusy {
		return nil, fmt.Errorf("%w: %s", ErrDateTimePeriodIsOccupied, dateTimePeriod)
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
		return nil, err
	}
	app := NewAppointmentAggregate(record, service, customer)
	if err := s.appointmentCreator.CreateAppointment(ctx, app); err != nil {
		return nil, err
	}
	if app.Id() == TemporalRecordId {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRecordId, app.Id())
	}
	return app, nil
}

func (s *SchedulingService) Schedule(
	ctx context.Context,
	now time.Time,
	preferredDate time.Time,
) (Schedule, error) {
	productionCalendar, err := s.productionCalendarLoader.ProductionCalendar(ctx)
	if err != nil {
		return Schedule{}, err
	}
	date := productionCalendar.DayOrNextWorkingDay(preferredDate)
	workingHours, err := s.workingHoursLoader.WorkingHours(ctx)
	if err != nil {
		return Schedule{}, err
	}
	dayTimePeriod, err := workingHours.ForDay(date).
		OmitPast(entity.GoTimeToDateTime(now)).
		ConsiderProductionCalendar(productionCalendar)
	if err != nil {
		return Schedule{}, err
	}
	busyPeriods, err := s.busyPeriodsLoader.BusyPeriods(ctx, date)
	if err != nil {
		return Schedule{}, err
	}
	workBreaks, err := s.workBreaksLoader.WorkBreaks(ctx)
	if err != nil {
		return Schedule{}, err
	}
	dayWorkBreaks, err := workBreaks.ForDay(date)
	if err != nil {
		return Schedule{}, err
	}
	schedulePeriods := NewSchedulePeriods(dayTimePeriod.Periods, busyPeriods, dayWorkBreaks)
	return NewSchedule(date, schedulePeriods, productionCalendar), nil
}
