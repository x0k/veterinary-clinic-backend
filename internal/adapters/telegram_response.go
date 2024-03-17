package adapters

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type TelegramResponseType int

const (
	TelegramText TelegramResponseType = iota
	TelegramQuery
	TelegramCallback
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

type TelegramCallbackResponse struct {
	Response *telebot.CallbackResponse
}

func (r TelegramCallbackResponse) Type() TelegramResponseType {
	return TelegramCallback
}
