package profiler

import (
	"net/http"
	"net/http/pprof"

	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
)

type Config struct {
	Enabled bool   `yaml:"enabled" env:"PROFILER_ENABLED"`
	Address string `yaml:"address" env:"PROFILER_ADDRESS"`
}

func UseHttpRouter(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return mux
}

func New(cfg *Config, log *logger.Logger) *module.Module {
	m := module.New(log.Logger, "profiler")
	if cfg.Enabled {
		m.Append(infra.NewHttpService(
			log,
			&http.Server{
				Addr:    cfg.Address,
				Handler: UseHttpRouter(http.NewServeMux()),
			},
		))
	}
	return m
}
