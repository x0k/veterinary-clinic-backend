package appointment_http_controller

import (
	"fmt"
	"net/http"
	"time"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	adapters_web_calendar "github.com/x0k/veterinary-clinic-backend/internal/adapters/web_calendar"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"gopkg.in/telebot.v3"
)

func UseWebCalendarRouter(
	mux *http.ServeMux,
	log *logger.Logger,
	bot *telebot.Bot,
	webCalendarAppOrigin adapters_web_calendar.AppOrigin,
	telegramIniDataParser *telegram_adapters.InitDataParser,
	scheduleUseCase *appointment_use_case.ScheduleUseCase[telegram_adapters.QueryResponse],
) *http.ServeMux {

	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", adapters_web_calendar.HandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webCalendarAppOrigin.String())
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s", adapters_web_calendar.HandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webCalendarAppOrigin.String())
		w.Header().Set("Vary", "Accept-Encoding, Origin")
		res, httpErr := httpx.JSONBody[adapters_web_calendar.AppResultResponse](log.Logger, jsonBodyDecoder, w, r)
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
		t, err := time.Parse(time.DateOnly, res.Data.SelectedDates[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		schedule, err := scheduleUseCase.Schedule(r.Context(), time.Now(), t)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to schedule",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, err = bot.AnswerWebApp(
			&telebot.Query{
				ID: data.QueryID,
			},
			schedule.Result,
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
	return mux
}
