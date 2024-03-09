package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

var (
	config_path string
)

func init() {
	flag.StringVar(&config_path, "config", os.Getenv("CONFIG_PATH"), "Config path")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	cfg := config.MustLoad(config_path)

	log := mustSetupLogger(&cfg.Logger)
	log.Info(ctx, "starting application", slog.String("log_level", cfg.Logger.Level))

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := start(ctx, wg, cfg, log); err != nil {
		log.Error(ctx, "failed to start application", sl.Err(err))
		cancel()
		wg.Wait()
		os.Exit(1)
	}

	log.Info(ctx, "application started")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info(ctx, "graceful shutdown")
	cancel()
	wg.Wait()
	log.Info(ctx, "application stopped")
}
