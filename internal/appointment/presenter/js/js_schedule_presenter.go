//go:build js && wasm

package appointment_js_presenter

import (
	"syscall/js"
	"time"

	"github.com/norunners/vert"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

func SchedulePresenter(
	now time.Time,
	schedule appointment.Schedule,
) (js.Value, error) {
	dto := appointment_js_adapters.ScheduleToDTO(schedule)
	return vert.ValueOf(dto).Value, nil
}
