package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type DetectChangesUseCase struct {
	log                        *logger.Logger
	appointmentChangesDetector *appointment.ChangesDetectorService
	publisher                  appointment.Publisher
}

func (u *DetectChangesUseCase) DetectChanges(ctx context.Context) {
	changes, err := u.appointmentChangesDetector.DetectChanges(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to detect changes", sl.Err(err))
		return
	}
	for _, change := range changes {
		if err := u.publisher.Publish(change); err != nil {
			u.log.Error(ctx, "failed to publish event", sl.Err(err))
		}
	}
}
