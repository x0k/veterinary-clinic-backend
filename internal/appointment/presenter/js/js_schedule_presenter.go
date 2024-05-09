//go:build js && wasm

package appointment_js_presenter

import (
	"time"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

func SchedulePresenter(
	now time.Time,
	schedule appointment.Schedule,
) (js_adapters.Result, error) {
	dto := appointment_js_adapters.ScheduleToDTO(schedule)
	return js_adapters.Ok(vert.ValueOf(dto)), nil
}
