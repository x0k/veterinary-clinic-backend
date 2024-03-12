package profiler_http_server

import (
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/config"
	"github.com/x0k/veterinary-clinic-backend/internal/controller/http/profiler"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type Server struct {
	shared.HttpService
}

const component_name = "profiler_http_server"

func New(cfg *config.ProfilerConfig, fataler shared.Fataler) *Server {
	mux := http.NewServeMux()
	profiler.UseRouter(mux)
	return &Server{
		HttpService: *shared.NewHttpService(
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
