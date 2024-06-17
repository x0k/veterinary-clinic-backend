//go:build js && wasm

package appointment_js_presenter

import (
	"syscall/js"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

func CustomerPresenter(
	customer appointment.CustomerEntity,
) (js_adapters.Result, error) {
	return js_adapters.Ok(
		js.ValueOf(
			customer.Id.String(),
		),
	), nil
}
