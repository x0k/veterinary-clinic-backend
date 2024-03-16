package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type CalendarDialogResult struct {
	Calendar struct {
		SelectedDates []string `json:"selectedDates"`
	} `json:"calendar"`
	WebAppInitData string `json:"webAppInitData"`
}

type HttpTelegramConfig struct {
	Token                    string
	CalendarInputHandlerPath string
	CalendarWebAppOrigin     string
}

func UseHttpTelegramRouter(
	log *logger.Logger,
	mux *http.ServeMux,
	clinicDialog *usecase.ClinicDialogUseCase[adapters.TelegramResponse],
	cfg *HttpTelegramConfig,
) {
	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}
	initDataParser := NewTelegramInitData(cfg.Token, time.Hour*24)

	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", cfg.CalendarInputHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", cfg.CalendarWebAppOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(fmt.Sprintf("POST %s", cfg.CalendarInputHandlerPath), func(w http.ResponseWriter, r *http.Request) {
		res, httpErr := httpx.JSONBody[CalendarDialogResult](log.Logger, jsonBodyDecoder, w, r)
		if httpErr != nil {
			http.Error(w, httpErr.Text, httpErr.Status)
			return
		}
		if len(res.Calendar.SelectedDates) == 0 {
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
				"failed to parse valid init data",
				slog.String("data", res.WebAppInitData),
				sl.Err(err),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		t, err := time.Parse(time.DateOnly, res.Calendar.SelectedDates[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		clinicDialog.FinishScheduleDialog(
			r.Context(),
			entity.Dialog{
				Id:     entity.DialogId(data.QueryID),
				UserId: entity.UserId(strconv.FormatInt(data.User.ID, 10)),
			},
			t,
		)
		w.WriteHeader(http.StatusOK)
	})
}
