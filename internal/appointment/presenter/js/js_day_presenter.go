//go:build js && wasm

package appointment_js_presenter

import (
	"syscall/js"
	"time"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
)

func DayPresenter(
	day time.Time,
) (js_adapters.Result, error) {
	return js_adapters.Ok(js.ValueOf(day.Format(time.RFC3339))), nil
}
