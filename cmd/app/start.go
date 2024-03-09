package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func start(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *logger.Logger) error {
	const op = "start"
	startProfiler(ctx, wg, log, cfg.Profiler)
	if err := startNewBot(ctx, cfg, wg, log); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
