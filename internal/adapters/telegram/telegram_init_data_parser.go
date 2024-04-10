package adapters_telegram

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
)

type InitDataParser struct {
	telegramToken adapters.TelegramToken
	expiredIn     time.Duration
}

func NewInitDataParser(telegramToken adapters.TelegramToken, expiredIn time.Duration) *InitDataParser {
	return &InitDataParser{
		telegramToken: telegramToken,
		expiredIn:     expiredIn,
	}
}

func (p *InitDataParser) Validate(data string) error {
	return initdata.Validate(data, string(p.telegramToken), p.expiredIn)
}

func (p *InitDataParser) Parse(data string) (initdata.InitData, error) {
	return initdata.Parse(data)
}
