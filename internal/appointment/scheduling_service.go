package appointment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrInvalidRecordId = errors.New("invalid record id")
var ErrDateTimePeriodIsOccupied = errors.New("date time period is occupied")
var ErrAnotherAppointmentIsAlreadyScheduled = errors.New("another appointment is already scheduled")
var ErrInvalidAppointmentStatusForCancel = errors.New("invalid appointment status")

type SchedulingService struct {
	log            *logger.Logger
	periodLocker   DateTimePeriodLocker
	periodUnLocker DateTimePeriodUnLocker

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
	periodLocker DateTimePeriodLocker,
	periodUnLocker DateTimePeriodUnLocker,
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
		periodLocker:                    periodLocker,
		periodUnLocker:                  periodUnLocker,
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

func (s *SchedulingService) MakeAppointment(
	ctx context.Context,
	now time.Time,
	appointmentDate time.Time,
	customer CustomerEntity,
	service ServiceEntity,
) (RecordEntity, error) {
	appointmentDateTime := shared.GoTimeToDateTime(appointmentDate)
	dateTimePeriod := shared.DateTimePeriod{
		Start: appointmentDateTime,
		End: shared.DateTime{
			Date: appointmentDateTime.Date,
			Time: shared.MakeTimeShifter(shared.Time{
				Minutes: service.DurationInMinutes.Int(),
			})(appointmentDateTime.Time),
		},
	}
	if err := s.periodLocker(ctx, dateTimePeriod); err != nil {
		return RecordEntity{}, err
	}
	defer func() {
		if err := s.periodUnLocker(ctx, dateTimePeriod); err != nil {
			s.log.Error(ctx, "failed to unlock period", sl.Err(err))
		}
	}()
	existedAppointment, err := s.customerActiveAppointmentLoader(ctx, customer.Id)
	if !errors.Is(err, shared.ErrNotFound) {
		if err != nil {
			return RecordEntity{}, err
		}
		return RecordEntity{}, fmt.Errorf("%w: %s", ErrAnotherAppointmentIsAlreadyScheduled, existedAppointment.Id)
	}
	productionCalendar, err := s.productionCalendar(ctx)
	if err != nil {
		return RecordEntity{}, err
	}
	busyPeriods, err := s.busyPeriodsLoader(ctx, appointmentDate)
	if err != nil {
		return RecordEntity{}, err
	}
	datWorkBreaks, err := s.dayWorkBreaks(ctx, appointmentDate)
	if err != nil {
		return RecordEntity{}, err
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
		return RecordEntity{}, err
	}
	if !freeTimeSlots.Includes(shared.TimePeriod{
		Start: dateTimePeriod.Start.Time,
		End:   dateTimePeriod.End.Time,
	}) {
		return RecordEntity{}, fmt.Errorf("%w: %s", ErrDateTimePeriodIsOccupied, dateTimePeriod)
	}
	title, err := RecordTitle(customer, service, now)
	if err != nil {
		return RecordEntity{}, err
	}
	record, err := NewRecord(
		TemporalRecordId,
		title,
		RecordAwaits,
		false,
		dateTimePeriod,
		customer.Id,
		service.Id,
		now,
	)
	if err != nil {
		return RecordEntity{}, err
	}
	if err := s.appointmentCreator(ctx, &record); err != nil {
		return RecordEntity{}, err
	}
	if record.Id == TemporalRecordId {
		return RecordEntity{}, fmt.Errorf("%w: %s", ErrInvalidRecordId, record.Id)
	}
	return record, nil
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
	busyPeriods, err := s.busyPeriodsLoader(ctx, appointmentDate)
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
	busyPeriods, err := s.busyPeriodsLoader(ctx, appointmentDate)
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
	customerId CustomerId,
) (RecordEntity, error) {
	rec, err := s.customerActiveAppointmentLoader(ctx, customerId)
	if err != nil {
		return RecordEntity{}, err
	}
	if rec.Status != RecordAwaits {
		return RecordEntity{}, fmt.Errorf("%w: %s", ErrInvalidAppointmentStatusForCancel, rec.Status)
	}
	return rec, s.appointmentRemover(ctx, rec.Id)
}

func (s *SchedulingService) productionCalendar(ctx context.Context) (ProductionCalendar, error) {
	pc, err := s.productionCalendarLoader(ctx)
	if err != nil {
		return ProductionCalendar{}, err
	}
	return pc.WithoutSaturdayWeekend(), nil
}

func (s *SchedulingService) dayWorkBreaks(ctx context.Context, day time.Time) (DayWorkBreaks, error) {
	workBreaks, err := s.workBreaksLoader(ctx)
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
	workingHours, err := s.workingHoursLoader(ctx)
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
