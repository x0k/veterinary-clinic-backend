//go:build js && wasm

package app_wasm

import (
	appointment_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
)

type LoggerConfig struct {
	Level int `js:"level"`
}

type Config struct {
	Logger      LoggerConfig                   `js:"logger"`
	Appointment appointment_wasm_module.Config `js:"appointment"`
}
