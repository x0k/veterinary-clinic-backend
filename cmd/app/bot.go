package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/bot"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func startNewBot(
	ctx context.Context,
	cfg *config.Config,
	wg *sync.WaitGroup,
	log *logger.Logger,
) error {
	const op = "startNewBot"
	b, err := bot.New(cfg, log)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	b.Start(ctx, wg)
	return nil
}
