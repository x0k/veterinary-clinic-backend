package appointment_http_controller

import (
	"context"
	"fmt"
	"net/http"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	web_calendar_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/web_calendar"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"gopkg.in/telebot.v3"
)

func useWebCalendarEndpoints(
	mux *http.ServeMux,
	log *logger.Logger,
	bot *telebot.Bot,
	endpointPath string,
	webCalendarAppOrigin web_calendar_adapters.AppOrigin,
	telegramInitDataParser telegram_adapters.InitDataParser,
	useCase func(context.Context, web_calendar_adapters.AppResultResponse) (telegram_adapters.QueryResponse, error),
) error {
	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", endpointPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webCalendarAppOrigin.String())
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(
		fmt.Sprintf("POST %s", endpointPath),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", webCalendarAppOrigin.String())
			w.Header().Set("Vary", "Accept-Encoding, Origin")
			res, httpErr := httpx.JSONBody[web_calendar_adapters.AppResultResponse](log.Logger, jsonBodyDecoder, w, r)
			if httpErr != nil {
				http.Error(w, httpErr.Text, httpErr.Status)
				return
			}
			if len(res.Data.SelectedDates) == 0 {
				log.Error(
					r.Context(),
					"no selected date",
				)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if err := telegramInitDataParser.Validate(res.WebAppInitData); err != nil {
				log.Error(
					r.Context(),
					"failed to validate init data",
					sl.Err(err),
				)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			data, err := telegramInitDataParser.Parse(res.WebAppInitData)
			if err != nil {
				log.Error(
					r.Context(),
					"failed to parse init data",
					sl.Err(err),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			result, err := useCase(r.Context(), res)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			_, err = bot.AnswerWebApp(
				&telebot.Query{
					ID: data.QueryID,
				},
				result.Result,
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
		},
	)
	return nil
}
