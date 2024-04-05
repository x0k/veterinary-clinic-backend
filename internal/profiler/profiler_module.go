package profiler_module

import (
	"net/http"
	"net/http/pprof"

	adapters_http "github.com/x0k/veterinary-clinic-backend/internal/adapters/http"
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
		srv := &http.Server{
			Addr:    cfg.Address,
			Handler: UseHttpRouter(http.NewServeMux()),
		}
		m.Append(adapters_http.NewService("profiler_http_server", srv, m))
	}
	return m
}
