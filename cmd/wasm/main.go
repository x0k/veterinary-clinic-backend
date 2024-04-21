//go:build js && wasm

package main

import (
	"errors"
	"log/slog"
	"syscall/js"

	"github.com/norunners/vert"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	app_wasm "github.com/x0k/veterinary-clinic-backend/internal/app/wasm"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

var ErrConfigExpected = errors.New("config argument expected")

func main() {
	js.Global().Set("__init_wasm", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js_adapters.RejectError(ErrConfigExpected)
		}
		cfgData := args[0]
		cfg := app_wasm.Config{}
		if err := vert.ValueOf(cfgData).AssignTo(&cfg); err != nil {
			return js_adapters.RejectError(err)
		}
		log := logger.New(
			slog.New(
				js_adapters.NewConsoleLoggerHandler(slog.Level(cfg.Logger.Level)),
			),
		)
		return app_wasm.New(&cfg, log)
	}))
	select {}
}
