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

type AppointmentService struct {
	log          *logger.Logger
	appointments AppointmentRepository
	clients      ClientRepository
	services     ServiceRepository
}

func NewAppointmentService(
	log *logger.Logger,
	appointments AppointmentRepository,
	clients ClientRepository,
	services ServiceRepository,
) *AppointmentService {
	return &AppointmentService{
		log:          log.With(slog.String("component", "appointment.AppointmentService")),
		appointments: appointments,
		clients:      clients,
		services:     services,
	}
}

func (s *AppointmentService) CreateAppointment(
	ctx context.Context,
	clientId ClientId,
	serviceId ServiceId,
	dateTimePeriod entity.DateTimePeriod,
) error {
	s.log.Debug(
		ctx, "create appointment",
		slog.String("client_id", clientId.String()),
		slog.String("service_id", serviceId.String()),
		slog.String("date_time_period", dateTimePeriod.String()),
	)
	client, err := s.clients.GetClient(ctx, clientId)
	if err != nil {
		return err
	}
	service, err := s.services.GetService(ctx, serviceId)
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
	appointment, err := NewAppointment(client, service, dateTimePeriod)
	if err != nil {
		return err
	}
	return s.appointments.SaveAppointment(ctx, appointment)
}
