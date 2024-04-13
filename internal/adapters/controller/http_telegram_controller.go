package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase/make_appointment"
)

func UseHttpTelegramRouter(
	mux *http.ServeMux,
	log *logger.Logger,
	query chan<- shared.DialogMessage[adapters.TelegramQueryResponse],
	calendarWebAppOrigin adapters.CalendarWebAppOrigin,
	telegramInitDataParser TelegramInitDataParser,
	schedule *usecase.ScheduleUseCase[adapters.TelegramQueryResponse],
	makeAppointmentDatePicker *make_appointment.DatePickerUseCase[adapters.TelegramQueryResponse],
) *http.ServeMux {
	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", adapters.CalendarWebHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", string(calendarWebAppOrigin))
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s", adapters.CalendarWebHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", string(calendarWebAppOrigin))
		w.Header().Set("Vary", "Accept-Encoding, Origin")
		res, httpErr := httpx.JSONBody[WebAppResultResponse](log.Logger, jsonBodyDecoder, w, r)
		if httpErr != nil {
			http.Error(w, httpErr.Text, httpErr.Status)
			return
		}
		if len(res.Data.SelectedDates) == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := telegramInitDataParser.Validate(res.WebAppInitData); err != nil {
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
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		t, err := time.Parse(time.DateOnly, res.Data.SelectedDates[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		schedule, err := schedule.Schedule(r.Context(), time.Now(), t)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to get schedule",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		query <- shared.DialogMessage[adapters.TelegramQueryResponse]{
			DialogId: shared.DialogId(data.QueryID),
			Message:  schedule,
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", adapters.MakeAppointmentDatePickerHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", string(calendarWebAppOrigin))
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s", adapters.MakeAppointmentDatePickerHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		res, httpErr := httpx.JSONBody[WebAppResultResponse](log.Logger, jsonBodyDecoder, w, r)
		if httpErr != nil {
			http.Error(w, httpErr.Text, httpErr.Status)
			return
		}
		if len(res.Data.SelectedDates) == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if res.State == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := telegramInitDataParser.Validate(res.WebAppInitData); err != nil {
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
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		selectedDate, err := time.Parse(time.DateOnly, res.Data.SelectedDates[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		datePicker, err := makeAppointmentDatePicker.DatePicker(r.Context(), shared.ServiceId(res.State), time.Now(), selectedDate)
		if err != nil {
			log.Error(
				r.Context(),
				"failed to get date picker",
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		query <- shared.DialogMessage[adapters.TelegramQueryResponse]{
			DialogId: shared.DialogId(data.QueryID),
			Message:  datePicker,
		}
		w.WriteHeader(http.StatusOK)
	})

	return mux
}
