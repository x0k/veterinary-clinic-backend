//go:build js && wasm

package app_wasm

import (
	appointment_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
)

type Config struct {
	Appointment appointment_wasm_module.Config `js:"appointment"`
}
