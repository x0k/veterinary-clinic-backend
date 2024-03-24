package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

type ArchiveRecordsUseCase struct {
	log         *logger.Logger
	archiveTime time.Time
	recordsRepo RecordsArchiver
}

func NewArchiveRecordsUseCase(log *logger.Logger, archiveTime time.Time, recordsRepo RecordsArchiver) *ArchiveRecordsUseCase {
	return &ArchiveRecordsUseCase{
		log:         log.With(slog.String("component", "usecase.ArchiveRecordsUseCase")),
		archiveTime: archiveTime,
		recordsRepo: recordsRepo,
	}
}

func (u *ArchiveRecordsUseCase) ArchiveRecords(ctx context.Context, now time.Time) error {
	targetTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		u.archiveTime.Hour(), u.archiveTime.Minute(), u.archiveTime.Second(),
		0, now.Location())
	if now.After(targetTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}
	timer := time.NewTimer(
		targetTime.Sub(now),
	)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil
	case <-timer.C:
		u.log.Info(ctx, "archiving records")
		return u.recordsRepo.ArchiveRecords(ctx)
	}
}
