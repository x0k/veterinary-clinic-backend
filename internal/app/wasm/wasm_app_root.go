//go:build js && wasm

package app_wasm

import (
	"context"
	"syscall/js"

	"github.com/jomei/notionapi"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	appointment_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	shared_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/shared/module/wasm"
)

func New(
	cfg *Config,
	log *logger.Logger,
) js.Value {
	ctx := context.Background()
	root := js_adapters.ObjectConstructor.New()
	sharedModule := shared_wasm_module.New()
	root.Set("shared", sharedModule)

	notion := notionapi.NewClient(cfg.Notion.Token)

	appointmentModule := appointment_wasm_module.New(
		ctx,
		&cfg.Appointment,
		log,
		notion,
	)
	root.Set("appointment", appointmentModule)
	return root
}
