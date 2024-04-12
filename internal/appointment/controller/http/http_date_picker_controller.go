package appointment_http_controller

import (
	"fmt"
	"net/http"
	"time"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_telegram_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"gopkg.in/telebot.v3"
)

func UseDatePickerRouter(
	mux *http.ServeMux,
	log *logger.Logger,
	bot *telebot.Bot,
	webCalendarAppOrigin web_calendar_adapters.AppOrigin,
	telegramIniDataParser *telegram_adapters.InitDataParser,
	appointmentDatePickerUseCase *appointment_telegram_use_case.AppointmentDatePickerUseCase[telegram_adapters.QueryResponse],
) error {

	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", web_calendar_adapters.DatePickerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webCalendarAppOrigin.String())
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s", web_calendar_adapters.DatePickerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webCalendarAppOrigin.String())
		w.Header().Set("Vary", "Accept-Encoding, Origin")
		res, httpErr := httpx.JSONBody[web_calendar_adapters.AppResultResponse](log.Logger, jsonBodyDecoder, w, r)
		if httpErr != nil {
			http.Error(w, httpErr.Text, httpErr.Status)
			return
		}
		if len(res.Data.SelectedDates) == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := telegramIniDataParser.Validate(res.WebAppInitData); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		data, err := telegramIniDataParser.Parse(res.WebAppInitData)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to parse init data",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		selectedDate, err := time.Parse(time.DateOnly, res.Data.SelectedDates[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		datePicker, err := appointmentDatePickerUseCase.DatePicker(
			r.Context(),
			appointment.NewServiceId(res.State),
			time.Now(),
			selectedDate,
		)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to get date picker",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, err = bot.AnswerWebApp(
			&telebot.Query{
				ID: data.QueryID,
			},
			datePicker.Result,
		)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to answer query",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return nil
}
