package adapters_cron

import (
	"context"
	"time"
)

type Task struct {
	name     string
	interval time.Duration
	task     func(context.Context, time.Time)
}

func NewTask(
	name string,
	interval time.Duration,
	task func(context.Context, time.Time),
) *Task {
	return &Task{
		name:     name,
		interval: interval,
		task:     task,
	}
}

func (c *Task) Name() string {
	return c.name
}

func (c *Task) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case now := <-ticker.C:
			c.task(ctx, now)
		}
	}
}
