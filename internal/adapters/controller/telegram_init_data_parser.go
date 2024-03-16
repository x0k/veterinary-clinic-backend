package controller

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type TelegramInitData struct {
	telegramToken string
	expiredIn     time.Duration
}

func NewTelegramInitData(telegramToken string, expiredIn time.Duration) *TelegramInitData {
	return &TelegramInitData{
		telegramToken: telegramToken,
		expiredIn:     expiredIn,
	}
}

func (p *TelegramInitData) Validate(data string) error {
	return initdata.Validate(data, p.telegramToken, p.expiredIn)
}

func (p *TelegramInitData) Parse(data string) (initdata.InitData, error) {
	return initdata.Parse(data)
}
