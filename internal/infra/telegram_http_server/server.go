package telegram_http_server

import (
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type Server struct {
	infra.HttpService
}

func New(
	log *logger.Logger,
	clinicSchedule *usecase.ClinicScheduleUseCase[adapters.TelegramQueryResponse],
	query chan<- entity.DialogMessage[adapters.TelegramQueryResponse],
	telegramHttpServerAddress infra.TelegramHttpServerAddress,
	telegramToken adapters.TelegramToken,
	calendarWebAppOrigin adapters.CalendarWebAppOrigin,
) *Server {
	mux := http.NewServeMux()
	controller.UseHttpTelegramRouter(
		log, mux,
		clinicSchedule,
		query,
		telegramToken,
		calendarWebAppOrigin,
	)
	return &Server{
		HttpService: *infra.NewHttpService(
			log,
			&http.Server{
				Addr:    string(telegramHttpServerAddress),
				Handler: infra.Logging(log, mux),
			},
		),
	}
}
