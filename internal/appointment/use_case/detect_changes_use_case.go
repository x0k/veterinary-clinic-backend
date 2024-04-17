package appointment_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const detectChangesUseCaseName = "appointment_use_case.DetectChangesUseCase"

type DetectChangesUseCase struct {
	log             *logger.Logger
	trackingService *appointment.TrackingService
	publisher       appointment.Publisher
}

func NewDetectChangesUseCase(
	log *logger.Logger,
	trackingService *appointment.TrackingService,
	publisher appointment.Publisher,
) *DetectChangesUseCase {
	return &DetectChangesUseCase{
		log:             log.With(sl.Component(detectChangesUseCaseName)),
		trackingService: trackingService,
		publisher:       publisher,
	}
}

func (u *DetectChangesUseCase) DetectChanges(ctx context.Context, now time.Time) {
	changes, err := u.trackingService.DetectChanges(ctx, now)
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
