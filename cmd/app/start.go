package main

import (
	"context"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func start(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *logger.Logger) error {
	const op = "start"
	startProfiler(ctx, wg, log, cfg.Profiler)
	return nil
}
