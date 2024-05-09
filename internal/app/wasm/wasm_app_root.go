//go:build js && wasm

package app_wasm

import (
	"syscall/js"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	appointment_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	shared_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/shared/module/wasm"
)

func New(
	cfg *Config,
	log *logger.Logger,
) js.Value {
	root := js_adapters.ObjectConstructor.New()
	sharedModule := shared_wasm_module.New()
	root.Set("shared", sharedModule)
	appointmentModule := appointment_wasm_module.New(
		&cfg.Appointment,
		log,
	)
	root.Set("appointment", appointmentModule)
	return root
}
