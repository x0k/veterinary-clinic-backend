package adapters_telegram

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type InitDataParser struct {
	telegramToken Token
	expiredIn     time.Duration
}

func NewInitDataParser(telegramToken Token, expiredIn time.Duration) *InitDataParser {
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
