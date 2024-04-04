package appointment_use_case

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

var ErrPeriodIsLocked = errors.New("periods is locked")
var ErrDateTimePeriodIsOccupied = errors.New("date time period is occupied")

type SchedulingService struct {
	log                       *logger.Logger
	periodsMu                 sync.Mutex
	periods                   []entity.DateTimePeriod
	appointmentPeriodsChecker AppointmentPeriodChecker
	appointmentCreator        AppointmentCreator
}

func NewSchedulingService(
	log *logger.Logger,
	appointmentPeriodChecker AppointmentPeriodChecker,
	appointmentCreator AppointmentCreator,
) *SchedulingService {
	return &SchedulingService{
		log:                       log.With(slog.String("component", "appointment.SchedulingService")),
		appointmentPeriodsChecker: appointmentPeriodChecker,
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
	customer appointment.CustomerEntity,
	service appointment.ServiceEntity,
	dateTimePeriod entity.DateTimePeriod,
) (*appointment.AppointmentAggregate, error) {
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
	record, err := appointment.NewRecord(
		appointment.NewRecordId(uuid.NewString()),
		appointment.RecordAwaits,
		false,
		dateTimePeriod,
		customer.Id,
		service.Id,
		now,
	)
	if err != nil {
		return nil, err
	}
	app := appointment.NewAppointmentAggregate(record, service, customer)
	if err := s.appointmentCreator.CreateAppointment(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}
