package profiler_http_server

import (
	"log/slog"
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/controller/http/profiler"
	"github.com/x0k/veterinary-clinic-backend/internal/infra/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

type Server struct {
	shared.HttpService
}

func New(cfg *config.ProfilerConfig, log *logger.Logger, fataler shared.Fataler) *Server {
	mux := http.NewServeMux()
	profiler.UseRouter(mux)
	return &Server{
		HttpService: *shared.NewHttpService(
			&http.Server{
				Addr:    cfg.Address,
				Handler: mux,
			},
			log.With(slog.String("component", "infra.profiler_http_server.Server")),
			fataler,
		),
	}
}

func (s *Server) Name() string {
	return "profiler_http_server"
}
