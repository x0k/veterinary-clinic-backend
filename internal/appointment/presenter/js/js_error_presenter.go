//go:build js && wasm

package appointment_js_presenter

import (
	"syscall/js"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
)

func ErrorPresenter(
	err error,
) (js.Value, error) {
	return js_adapters.RejectError(err), nil
}
