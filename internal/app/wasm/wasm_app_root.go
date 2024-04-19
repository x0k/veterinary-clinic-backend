//go:build js && wasm

package app_wasm

import (
	"syscall/js"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment/module/appointment_wasm_module"
)

func New(
	cfgData js.Value,
) js.Value {
	cfg := Config{}
	if err := vert.ValueOf(cfgData).AssignTo(&cfg); err != nil {
		return js_adapters.RejectError(err)
	}
	appointmentModule := appointment_wasm_module.New(
		&cfg.Appointment,
	)
}
