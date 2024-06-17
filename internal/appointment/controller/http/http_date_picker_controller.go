package appointment_http_controller

import (
	"context"
	"net/http"
	"time"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

func UseDatePickerRouter(
	mux *http.ServeMux,
	log *logger.Logger,
	bot *telebot.Bot,
	webCalendarAppOrigin web_calendar_adapters.AppOrigin,
	telegramIniDataParser telegram_adapters.InitDataParser,
	appointmentDatePickerUseCase *appointment_telegram_use_case.AppointmentDatePickerUseCase[telegram_adapters.QueryResponse],
) error {
	return useWebCalendarEndpoints(
		mux, log, bot,
		web_calendar_adapters.DatePickerPath,
		webCalendarAppOrigin,
		telegramIniDataParser,
		func(ctx context.Context, res web_calendar_adapters.AppResultResponse) (telegram_adapters.QueryResponse, error) {
			selectedDate, err := time.Parse(time.DateOnly, res.Data.SelectedDates[0])
			if err != nil {
				log.Error(
					ctx,
					"failed to parse selected date",
					sl.Err(err),
				)
				return telegram_adapters.QueryResponse{}, err
			}
			now := shared.NewUTCTime(time.Now())
			utcDate := shared.NewUTCTime(selectedDate)
			datePicker, err := appointmentDatePickerUseCase.DatePicker(
				ctx,
				appointment.NewServiceId(res.State),
				now,
				utcDate,
			)
			if err != nil {
				log.Error(
					ctx,
					"failed to get date picker",
					sl.Err(err),
				)
				return telegram_adapters.QueryResponse{}, err
			}
			return datePicker, nil
		},
	)
}
