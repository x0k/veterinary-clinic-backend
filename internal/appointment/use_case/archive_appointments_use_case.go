package appointment_use_case

import (
	"context"
	"fmt"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const archiveAppointmentsUseCaseName = "appointment_use_case.ArchiveAppointmentsUseCase"

type ArchiveAppointmentsUseCase struct {
	log             *logger.Logger
	archiveHour     int
	archiveMinutes  int
	recordsArchiver appointment.RecordsArchiver
}

func NewArchiveAppointmentsUseCase(
	log *logger.Logger,
	archiveHour int,
	archiveMinutes int,
	recordsArchiver appointment.RecordsArchiver,
) *ArchiveAppointmentsUseCase {
	return &ArchiveAppointmentsUseCase{
		log:             log.With(sl.Component(archiveAppointmentsUseCaseName)),
		archiveHour:     archiveHour,
		archiveMinutes:  archiveMinutes,
		recordsArchiver: recordsArchiver,
	}
}

func (u *ArchiveAppointmentsUseCase) ArchiveRecords(ctx context.Context, now time.Time) {
	targetTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		u.archiveHour, u.archiveMinutes,
		0, 0, now.Location())
	if now.After(targetTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}
	diff := targetTime.Sub(now)
	fmt.Println("Now ", now)
	fmt.Println("Target ", targetTime)
	fmt.Println("Wait ", diff.Seconds())
	timer := time.NewTimer(diff)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		u.log.Info(ctx, "archiving records")
		if err := u.recordsArchiver(ctx); err != nil {
			u.log.Error(ctx, "failed to archive records", sl.Err(err))
		}
	}
}
