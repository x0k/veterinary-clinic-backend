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
	CalendarInputHandlerPath string
	CalendarWebAppOrigin     string
	Address                  string
}

func New(
	log *logger.Logger,
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
	initDataValidator controller.TelegramInitDataParser,
	cfg *Config,
) *Server {
	mux := http.NewServeMux()
	controller.UseHttpTelegramRouter(
		log, mux,
		clinicDialog,
		initDataValidator,
		&controller.HttpTelegramConfig{
			CalendarInputHandlerPath: cfg.CalendarInputHandlerPath,
			CalendarWebAppOrigin:     cfg.CalendarWebAppOrigin,
		},
	)
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
