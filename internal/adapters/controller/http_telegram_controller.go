package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type WebAppResultResponse struct {
	Data struct {
		SelectedDates []string `json:"selectedDates"`
	} `json:"data"`
	WebAppInitData string `json:"webAppInitData"`
}

func UseHttpTelegramRouter(
	log *logger.Logger,
	mux *http.ServeMux,
	clinicSchedule *usecase.ClinicScheduleUseCase[adapters.TelegramQueryResponse],
	query chan<- entity.DialogMessage[adapters.TelegramQueryResponse],
	telegramToken adapters.TelegramToken,
	calendarWebAppOrigin adapters.CalendarWebAppOrigin,
) {
	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}
	initDataParser := NewTelegramInitData(telegramToken, time.Hour*24)

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", adapters.CalendarWebHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", string(calendarWebAppOrigin))
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s", adapters.CalendarWebHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		res, httpErr := httpx.JSONBody[WebAppResultResponse](log.Logger, jsonBodyDecoder, w, r)
		if httpErr != nil {
			http.Error(w, httpErr.Text, httpErr.Status)
			return
		}
		if len(res.Data.SelectedDates) == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := initDataParser.Validate(res.WebAppInitData); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		data, err := initDataParser.Parse(res.WebAppInitData)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to parse init data",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		t, err := time.Parse(time.DateOnly, res.Data.SelectedDates[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		schedule, err := clinicSchedule.Schedule(r.Context(), time.Now(), t)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to get schedule",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		query <- entity.DialogMessage[adapters.TelegramQueryResponse]{
			DialogId: entity.DialogId(data.QueryID),
			Message:  schedule,
		}
		w.WriteHeader(http.StatusOK)
	})
}
