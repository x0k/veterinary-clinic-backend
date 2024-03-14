package shared

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type TelegramResponseType int

const (
	TelegramText TelegramResponseType = iota
	TelegramQuery
)

type TelegramResponse interface {
	Type() TelegramResponseType
}

type TelegramTextResponse struct {
	Text    string
	Options *telebot.SendOptions
}

func (r TelegramTextResponse) Type() TelegramResponseType {
	return TelegramText
}

var escapeRegExp = regexp.MustCompile(`([.!-])`)

func EscapeTelegramMarkdownString(text string) string {
	return escapeRegExp.ReplaceAllString(text, "\\$1")
}

type TelegramQueryResponse struct {
	Result telebot.Result
}

func (r TelegramQueryResponse) Type() TelegramResponseType {
	return TelegramQuery
}
