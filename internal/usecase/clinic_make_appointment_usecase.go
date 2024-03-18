package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type clinicMakeAppointmentPresenter[R any] interface {
	RenderServicesList(services []entity.Service) (R, error)
	RenderDatePicker(serviceId entity.ServiceId, schedule entity.Schedule) (R, error)
	RenderTimePicker(serviceId entity.ServiceId, appointmentDate time.Time, slots entity.SampledFreeTimeSlots) (R, error)
	RenderAppointmentConfirmation(service entity.Service, appointmentDate time.Time) (R, error)
	RenderAppointmentInfo(appointment entity.Record) (R, error)
}

type makeAppointmentClinicRecordsRepo interface {
	clinicRecordsCreator
	clinicRecordsChecker
}

type makeAppointmentClinicServicesRepo interface {
	clinicServicesLoader
	clinicServiceLoader
}

type ClinicMakeAppointmentUseCase[R any] struct {
	sampleRateInMinutes    entity.SampleRateInMinutes
	recordsRepo            makeAppointmentClinicRecordsRepo
	clinicServicesRepo     makeAppointmentClinicServicesRepo
	productionCalendarRepo productionCalendarRepo
	openingHoursRepo       openingHoursRepo
	busyPeriodsRepo        busyPeriodsRepo
	workBreaksRepo         workBreaksRepo
	presenter              clinicMakeAppointmentPresenter[R]
}

func NewClinicMakeAppointmentUseCase[R any](
	sampleRateInMinutes entity.SampleRateInMinutes,
	recordsRepo makeAppointmentClinicRecordsRepo,
	clinicServicesRepo makeAppointmentClinicServicesRepo,
	productionCalendarRepo productionCalendarRepo,
	openingHoursRepo openingHoursRepo,
	busyPeriodsRepo busyPeriodsRepo,
	workBreaksRepo workBreaksRepo,
	presenter clinicMakeAppointmentPresenter[R],
) *ClinicMakeAppointmentUseCase[R] {
	return &ClinicMakeAppointmentUseCase[R]{
		sampleRateInMinutes:    sampleRateInMinutes,
		recordsRepo:            recordsRepo,
		clinicServicesRepo:     clinicServicesRepo,
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
		workBreaksRepo:         workBreaksRepo,
		presenter:              presenter,
	}
}

func (u *ClinicMakeAppointmentUseCase[R]) CheckExistingAppointment(ctx context.Context, userId entity.UserId) (bool, error) {
	return u.recordsRepo.Exists(ctx, userId)
}

func (u *ClinicMakeAppointmentUseCase[R]) Services(ctx context.Context) (R, error) {
	services, err := u.clinicServicesRepo.Services(ctx)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderServicesList(services)
}

func (u *ClinicMakeAppointmentUseCase[R]) DatePicker(
	ctx context.Context,
	serviceId entity.ServiceId,
	now time.Time,
) (R, error) {
	schedule, err := fetchAndCalculateSchedule(
		ctx,
		now,
		now,
		u.productionCalendarRepo,
		u.openingHoursRepo,
		u.busyPeriodsRepo,
		u.workBreaksRepo,
	)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderDatePicker(serviceId, schedule)
}

func (u *ClinicMakeAppointmentUseCase[R]) TimePicker(
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

func (u *ClinicMakeAppointmentUseCase[R]) AppointmentConfirmation(
	ctx context.Context,
	serviceId entity.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	service, err := u.clinicServicesRepo.Load(ctx, serviceId)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderAppointmentConfirmation(service, appointmentDateTime)
}

func (u *ClinicMakeAppointmentUseCase[R]) MakeAppointment(
	ctx context.Context,
	user entity.User,
	serviceId entity.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	record, err := u.recordsRepo.Create(
		ctx,
		user,
		serviceId,
		appointmentDateTime,
	)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderAppointmentInfo(record)
}
