package profiler_http_server

import (
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

type Server struct {
	infra.HttpService
}

func New(log *logger.Logger, cfg *config.ProfilerConfig) *Server {
	mux := http.NewServeMux()
	controller.UseHttpProfilerRouter(mux)
	return &Server{
		HttpService: *infra.NewHttpService(
			log,
			&http.Server{
				Addr:    cfg.Address,
				Handler: mux,
			},
		),
	}
}
