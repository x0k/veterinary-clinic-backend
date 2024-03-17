package controller

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
)

type TelegramInitData struct {
	telegramToken adapters.TelegramToken
	expiredIn     time.Duration
}

func NewTelegramInitData(telegramToken adapters.TelegramToken, expiredIn time.Duration) *TelegramInitData {
	return &TelegramInitData{
		telegramToken: telegramToken,
		expiredIn:     expiredIn,
	}
}

func (p *TelegramInitData) Validate(data string) error {
	return initdata.Validate(data, string(p.telegramToken), p.expiredIn)
}

func (p *TelegramInitData) Parse(data string) (initdata.InitData, error) {
	return initdata.Parse(data)
}
