//go:build js && wasm

package appointment_js_presenter

import (
	"syscall/js"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
)

func OkPresenter() (js_adapters.Result, error) {
	return js_adapters.Ok(js.ValueOf(true)), nil
}
