package telegram_adapters

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type InitDataParser interface {
	Validate(data string) error
	Parse(data string) (initdata.InitData, error)
}

type initDataParser struct {
	telegramToken Token
	expiredIn     time.Duration
}

func NewInitDataParser(telegramToken Token, expiredIn time.Duration) InitDataParser {
	return &initDataParser{
		telegramToken: telegramToken,
		expiredIn:     expiredIn,
	}
}

func (p *initDataParser) Validate(data string) error {
	return initdata.Validate(data, string(p.telegramToken), p.expiredIn)
}

func (p *initDataParser) Parse(data string) (initdata.InitData, error) {
	return initdata.Parse(data)
}
