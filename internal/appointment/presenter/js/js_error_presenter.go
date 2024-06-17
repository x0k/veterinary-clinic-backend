//go:build js && wasm

package appointment_js_presenter

import (
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
)

func ErrorPresenter(
	err error,
) (js_adapters.Result, error) {
	return js_adapters.Fail(err), nil
}
