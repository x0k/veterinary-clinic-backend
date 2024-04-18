package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const updateAppointmentsStateUseCaseName = "appointment_use_case.UpdateAppointmentsStateUseCase"

type UpdateAppointmentsStateUseCase struct {
	log             *logger.Logger
	trackingService *appointment.TrackingService
}

func NewUpdateAppointmentsStateUseCase(
	log *logger.Logger,
	trackingService *appointment.TrackingService,
) *UpdateAppointmentsStateUseCase {
	return &UpdateAppointmentsStateUseCase{
		log:             log.With(sl.Component(updateAppointmentsStateUseCaseName)),
		trackingService: trackingService,
	}
}

func (u *UpdateAppointmentsStateUseCase) AddAppointment(ctx context.Context, app appointment.RecordEntity) {
	if err := u.trackingService.AddAppointment(ctx, app); err != nil {
		u.log.Error(ctx, "failed to add appointment", sl.Err(err))
	}
}

func (u *UpdateAppointmentsStateUseCase) RemoveAppointment(ctx context.Context, app appointment.RecordEntity) {
	if err := u.trackingService.RemoveAppointment(ctx, app); err != nil {
		u.log.Error(ctx, "failed to remove appointment", sl.Err(err))
	}
}
