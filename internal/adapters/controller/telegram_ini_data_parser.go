package controller

import initdata "github.com/telegram-mini-apps/init-data-golang"

type TelegramInitDataParser interface {
	Validate(data string) error
	Parse(data string) (initdata.InitData, error)
}
