//go:build js && wasm

package app_wasm

import (
	"log/slog"
	"syscall/js"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	appointment_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func New(
	cfgData js.Value,
) js.Value {
	cfg := Config{}
	if err := vert.ValueOf(cfgData).AssignTo(&cfg); err != nil {
		return js_adapters.RejectError(err)
	}
	log := logger.New(
		slog.New(
			js_adapters.NewConsoleLoggerHandler(slog.Level(cfg.Logger.Level)),
		),
	)

	root := js_adapters.ObjectConstructor.New()
	appointmentModule, err := appointment_wasm_module.New(
		&cfg.Appointment,
		log,
	)
	if err != nil {
		return js_adapters.RejectError(err)
	}
	root.Set("appointment", appointmentModule)
	return root
}
