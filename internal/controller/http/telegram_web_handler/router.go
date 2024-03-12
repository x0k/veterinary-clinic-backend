package telegram_web_handler

import (
	"fmt"
	"net/http"

	initdata "github.com/telegram-mini-apps/init-data-golang"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/httpx"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type CalendarDialogResult struct {
	Calendar struct {
		SelectedDates []string `json:"selectedDates"`
	} `json:"calendar"`
	WebAppInitData string `json:"webAppInitData"`
}

type TelegramInitDataParser interface {
	Validate(data string) error
	Parse(data string) (initdata.InitData, error)
}

func UseRouter(
	log *logger.Logger,
	mux *http.ServeMux,
	clinicDialog *usecase.ClinicDialogUseCase[shared.TelegramResponse],
	initDataParser TelegramInitDataParser,
) {
	jsonBodyDecoder := &httpx.JsonBodyDecoder{
		MaxBytes: 1 * 1024 * 1024,
	}
	mux.HandleFunc("POST /calendar-input-result", func(w http.ResponseWriter, r *http.Request) {
		res, err := httpx.JSONBody[CalendarDialogResult](log.Logger, jsonBodyDecoder, w, r)
		if err != nil {
			http.Error(w, err.Text, err.Status)
			return
		}
		if err := initDataParser.Validate(res.WebAppInitData); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		data, err2 := initDataParser.Parse(res.WebAppInitData)
		if err2 != nil {
			log.Error(r.Context(), "failed to parse init data", sl.Err(err2))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Println(data)
		fmt.Println(res.Calendar)
	})
}
