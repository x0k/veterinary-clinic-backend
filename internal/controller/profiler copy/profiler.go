package profiler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/pprof"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

type Profiler struct {
	enabled bool
	log     *logger.Logger
}

func New(cfg *config.ProfilerConfig, log *logger.Logger) *Profiler {
	return &Profiler{
		enabled: cfg.Enabled,
		log:     log.With(slog.String("component", "controller.profiler")),
	}
}

func (p *Profiler) Start(ctx context.Context) {
	if !p.enabled {
		p.log.Info(ctx, "disabled")
		return
	}
	p.log.Info(ctx, "starting", slog.String("address", p.address))
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	p.srv = &http.Server{
		Addr:    p.address,
		Handler: mux,
	}
	go func() {
		if err := p.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			p.log.Error(ctx, "failed to start", sl.Err(err))
		}
	}()
}

func (p *Profiler) Stop(ctx context.Context) {
	if !p.enabled {
		return
	}
	p.log.Info(ctx, "stopping")
	if err := p.srv.Shutdown(ctx); err != nil {
		p.log.Error(ctx, "failed to stop", sl.Err(err))
	}
}
