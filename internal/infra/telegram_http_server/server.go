package telegram_http_server

import (
	"net/http"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters/controller"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/infra"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase/clinic_make_appointment"
)

type Server struct {
	infra.HttpService
}

func New(
	log *logger.Logger,
	query chan<- entity.DialogMessage[adapters.TelegramQueryResponse],
	telegramHttpServerAddress infra.TelegramHttpServerAddress,
	calendarWebAppOrigin adapters.CalendarWebAppOrigin,
	clinicSchedule *usecase.ClinicScheduleUseCase[adapters.TelegramQueryResponse],
	telegramInitDataParser controller.TelegramInitDataParser,
	makeAppointmentDatePicker *clinic_make_appointment.DatePickerUseCase[adapters.TelegramQueryResponse],
) *Server {
	mux := http.NewServeMux()
	controller.UseHttpTelegramRouter(
		mux,
		log,
		query,
		calendarWebAppOrigin,
		telegramInitDataParser,
		clinicSchedule,
		makeAppointmentDatePicker,
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
