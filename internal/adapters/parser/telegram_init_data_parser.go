package telegram_init_data_parser

import (
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type Parser struct {
	telegramToken string
	expiredIn     time.Duration
}

func New(telegramToken string, expiredIn time.Duration) *Parser {
	return &Parser{
		telegramToken: telegramToken,
		expiredIn:     expiredIn,
	}
}

func (p *Parser) Validate(data string) error {
	return initdata.Validate(data, p.telegramToken, p.expiredIn)
}

func (p *Parser) Parse(data string) (initdata.InitData, error) {
	return initdata.Parse(data)
}
