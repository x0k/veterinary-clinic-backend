package appointment

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

type MakeAppointmentUseCase struct {
	scheduling *SchedulingService
	customers  CustomerRepository
	services   ServiceRepository
}

func NewMakeAppointmentUseCase(
	log *logger.Logger,
	appointments *SchedulingService,
	customers CustomerRepository,
	services ServiceRepository,
) *MakeAppointmentUseCase {
	return &MakeAppointmentUseCase{
		scheduling: appointments,
		customers:  customers,
		services:   services,
	}
}

func (s *MakeAppointmentUseCase) CreateAppointment(
	ctx context.Context,
	customerId CustomerId,
	serviceId ServiceId,
	dateTimePeriod entity.DateTimePeriod,
) error {
	customer, err := s.customers.Customer(ctx, customerId)
	if err != nil {
		return err
	}
	service, err := s.services.Service(ctx, serviceId)
	if err != nil {
		return err
	}
	return s.scheduling.MakeAppointment(ctx, customer, service, dateTimePeriod)
}
