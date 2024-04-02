package appointment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

var ErrDateTimePeriodIsOccupied = errors.New("date time period is occupied")

type MakeAppointmentUseCase struct {
	log          *logger.Logger
	appointments AppointmentRepository
	customers    CustomerRepository
	services     ServiceRepository
}

func NewMakeAppointmentUseCase(
	log *logger.Logger,
	appointments AppointmentRepository,
	customers CustomerRepository,
	services ServiceRepository,
) *MakeAppointmentUseCase {
	return &MakeAppointmentUseCase{
		log:          log.With(slog.String("component", "appointment.AppointmentService")),
		appointments: appointments,
		customers:    customers,
		services:     services,
	}
}

func (s *MakeAppointmentUseCase) CreateAppointment(
	ctx context.Context,
	customerId CustomerId,
	serviceId ServiceId,
	dateTimePeriod entity.DateTimePeriod,
) error {
	s.log.Debug(
		ctx, "create appointment",
		slog.String("customer_id", customerId.String()),
		slog.String("service_id", serviceId.String()),
		slog.String("date_time_period", dateTimePeriod.String()),
	)
	customer, err := s.customers.Customer(ctx, customerId)
	if err != nil {
		return err
	}
	service, err := s.services.Service(ctx, serviceId)
	if err != nil {
		return err
	}
	if err := s.appointments.LockPeriod(ctx, dateTimePeriod); err != nil {
		return err
	}
	defer func() {
		if err := s.appointments.UnLockPeriod(ctx, dateTimePeriod); err != nil {
			s.log.Error(ctx, "failed to release appointments lock", sl.Err(err))
		}
	}()
	isBusy, err := s.appointments.IsAppointmentPeriodBusy(ctx, dateTimePeriod)
	if err != nil {
		return err
	}
	if isBusy {
		return fmt.Errorf("%w: %s", ErrDateTimePeriodIsOccupied, dateTimePeriod)
	}
	appointment, err := NewAppointment(customer, service, dateTimePeriod)
	if err != nil {
		return err
	}
	return s.appointments.SaveAppointment(ctx, appointment)
}
