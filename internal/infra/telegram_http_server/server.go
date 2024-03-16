package telegram_http_server

import (
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type Server struct {
	infra.HttpService
}

type Config struct {
	Token                    string
	CalendarInputHandlerPath string
	CalendarWebAppOrigin     string
	Address                  string
}

func New(
	log *logger.Logger,
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
	cfg *Config,
) *Server {
	mux := http.NewServeMux()
	controller.UseHttpTelegramRouter(
		log, mux,
		clinicDialog,
		&controller.HttpTelegramConfig{
			Token:                    cfg.Token,
			CalendarInputHandlerPath: cfg.CalendarInputHandlerPath,
			CalendarWebAppOrigin:     cfg.CalendarWebAppOrigin,
		},
	)
	infra.Logging(log, mux)
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
