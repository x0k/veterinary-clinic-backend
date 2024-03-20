package controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type Cron struct {
	log      *logger.Logger
	interval time.Duration
	task     func(context.Context, time.Time) error
}

func NewCron(
	log *logger.Logger,
	interval time.Duration,
	task func(context.Context, time.Time) error,
) *Cron {
	return &Cron{
		log:      log.With(slog.String("component", "adapters.controller.Cron")),
		interval: interval,
		task:     task,
	}
}

func (c *Cron) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case now := <-ticker.C:
			if err := c.task(ctx, now); err != nil {
				c.log.Error(ctx, "failed to run cron task", sl.Err(err))
			}
		}
	}
}
