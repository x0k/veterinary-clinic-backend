package clinic_make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type timeSlotPickerPresenter[R any] interface {
	RenderTimePicker(serviceId entity.ServiceId, appointmentDate time.Time, slots entity.SampledFreeTimeSlots) (R, error)
}

type TimeSlotPickerUseCase[R any] struct {
	sampleRateInMinutes    entity.SampleRateInMinutes
	productionCalendarRepo usecase.ProductionCalendarLoader
	openingHoursRepo       usecase.OpeningHoursLoader
	busyPeriodsRepo        usecase.BusyPeriodsLoader
	workBreaksRepo         usecase.WorkBreaksLoader
	clinicServicesRepo     usecase.ClinicServiceLoader
	presenter              timeSlotPickerPresenter[R]
}

func NewTimeSlotPickerUseCase[R any](
	sampleRateInMinutes entity.SampleRateInMinutes,
	productionCalendarRepo usecase.ProductionCalendarLoader,
	openingHoursRepo usecase.OpeningHoursLoader,
	busyPeriodsRepo usecase.BusyPeriodsLoader,
	workBreaksRepo usecase.WorkBreaksLoader,
	clinicServicesRepo usecase.ClinicServiceLoader,
	presenter timeSlotPickerPresenter[R],
) *TimeSlotPickerUseCase[R] {
	return &TimeSlotPickerUseCase[R]{
		sampleRateInMinutes:    sampleRateInMinutes,
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
		workBreaksRepo:         workBreaksRepo,
		clinicServicesRepo:     clinicServicesRepo,
		presenter:              presenter,
	}
}

func (u *TimeSlotPickerUseCase[R]) TimePicker(
	ctx context.Context,
	serviceId entity.ServiceId,
	now time.Time,
	appointmentDate time.Time,
) (R, error) {
	productionCalendar, err := u.productionCalendarRepo.ProductionCalendar(ctx)
	if err != nil {
		return *new(R), err
	}
	openingHours, err := u.openingHoursRepo.OpeningHours(ctx)
	if err != nil {
		return *new(R), err
	}
	freePeriods, err := entity.CalculateFreePeriods(
		productionCalendar,
		openingHours,
		now,
		appointmentDate,
	)
	if err != nil {
		return *new(R), err
	}
	busyPeriods, err := u.busyPeriodsRepo.BusyPeriods(ctx, appointmentDate)
	if err != nil {
		return *new(R), err
	}
	allWorkBreaks, err := u.workBreaksRepo.WorkBreaks(ctx)
	if err != nil {
		return *new(R), err
	}
	workBreaks, err := entity.CalculateWorkBreaks(allWorkBreaks, appointmentDate)
	if err != nil {
		return *new(R), err
	}
	freeTimeSlots := entity.CalculateFreeTimeSlots(
		freePeriods,
		busyPeriods,
		workBreaks,
	)
	clinicService, err := u.clinicServicesRepo.Load(ctx, serviceId)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderTimePicker(serviceId, appointmentDate, entity.SampleFreeTimeSlots(
		clinicService.DurationInMinutes,
		u.sampleRateInMinutes,
		freeTimeSlots,
	))
}
