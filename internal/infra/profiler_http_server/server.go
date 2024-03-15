package profiler_http_server

import (
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
)

type Server struct {
	infra.HttpService
}

const component_name = "profiler_http_server"

func New(cfg *config.ProfilerConfig, fataler infra.Fataler) *Server {
	mux := http.NewServeMux()
	controller.UseHttpProfilerRouter(mux)
	return &Server{
		HttpService: *infra.NewHttpService(
			component_name,
			&http.Server{
				Addr:    cfg.Address,
				Handler: mux,
			},
			fataler,
		),
	}
}

func (s *Server) Name() string {
	return component_name
}
