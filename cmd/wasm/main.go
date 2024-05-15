//go:build js && wasm

package main

import (
	"errors"
	"log/slog"
	"syscall/js"

	"github.com/x0k/vert"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	app_wasm "github.com/x0k/veterinary-clinic-backend/internal/app/wasm"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

var ErrConfigExpected = errors.New("config argument expected")

func main() {
	js.Global().Set("initRootDomain", js_adapters.Sync(func(args []js.Value) js_adapters.Result {
		if len(args) < 1 {
			return js_adapters.Fail(ErrConfigExpected)
		}
		cfgData := args[0]
		cfg := app_wasm.Config{}
		if err := vert.Assign(cfgData, &cfg); err != nil {
			return js_adapters.Fail(err)
		}
		log := logger.New(
			slog.New(
				js_adapters.NewConsoleLoggerHandler(slog.Level(cfg.Logger.Level)),
			),
		)
		return js_adapters.Ok(app_wasm.New(&cfg, log))
	}))
	select {}
}
