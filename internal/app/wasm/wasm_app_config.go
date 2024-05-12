//go:build js && wasm

package app_wasm

import (
	"github.com/jomei/notionapi"
	appointment_wasm_module "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
)

type LoggerConfig struct {
	Level int `js:"level"`
}

type NotionConfig struct {
	Token notionapi.Token `js:"token"`
}

type Config struct {
	Logger      LoggerConfig                   `js:"logger"`
	Notion      NotionConfig                   `js:"notion"`
	Appointment appointment_wasm_module.Config `js:"appointment"`
}
