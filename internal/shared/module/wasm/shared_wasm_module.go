//go:build js && wasm

package shared_wasm_module

import (
	"syscall/js"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	shared_js_controller "github.com/x0k/veterinary-clinic-backend/internal/shared/controller/js"
)

func New() js.Value {
	m := js_adapters.ObjectConstructor.New()
	// DateTime controller
	shared_js_controller.NewDateTime(m)
	return m
}
