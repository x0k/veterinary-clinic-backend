package adapters_telegram

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
)

type InitData struct {
	telegramToken adapters.TelegramToken
	expiredIn     time.Duration
}

func NewInitData(telegramToken adapters.TelegramToken, expiredIn time.Duration) *InitData {
	return &InitData{
		telegramToken: telegramToken,
		expiredIn:     expiredIn,
	}
}

func (p *InitData) Validate(data string) error {
	return initdata.Validate(data, string(p.telegramToken), p.expiredIn)
}

func (p *InitData) Parse(data string) (initdata.InitData, error) {
	return initdata.Parse(data)
}
