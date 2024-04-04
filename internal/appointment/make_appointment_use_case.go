package appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type MakeAppointmentUseCase struct {
	scheduling *SchedulingService
	customers  CustomerLoader
	services   ServiceLoader
}

func NewMakeAppointmentUseCase(
	appointments *SchedulingService,
	customers CustomerLoader,
	services ServiceLoader,
) *MakeAppointmentUseCase {
	return &MakeAppointmentUseCase{
		scheduling: appointments,
		customers:  customers,
		services:   services,
	}
}

func (s *MakeAppointmentUseCase) CreateAppointment(
	ctx context.Context,
	now time.Time,
	customerId CustomerId,
	serviceId ServiceId,
	dateTimePeriod entity.DateTimePeriod,
) (*AppointmentAggregate, error) {
	customer, err := s.customers.Customer(ctx, customerId)
	if err != nil {
		return nil, err
	}
	service, err := s.services.Service(ctx, serviceId)
	if err != nil {
		return nil, err
	}
	return s.scheduling.MakeAppointment(ctx, now, customer, service, dateTimePeriod)
}
