//go:build js && wasm

package app_wasm

import (
	"context"
	"net/http"
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

	httpClient := &http.Client{}

	notion := notionapi.NewClient(
		cfg.Notion.Token,
		notionapi.WithHTTPClient(httpClient),
	)

	appointmentModule := appointment_wasm_module.New(
		ctx,
		&cfg.Appointment,
		log,
		httpClient,
		notion,
	)
	root.Set("appointment", appointmentModule)
	return root
}
