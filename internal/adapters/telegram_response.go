package adapters

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type TelegramTextResponse struct {
	Text    string
	Options *telebot.SendOptions
}

var escapeRegExp = regexp.MustCompile(`([.!-])`)

func EscapeTelegramMarkdownString(text string) string {
	return escapeRegExp.ReplaceAllString(text, "\\$1")
}

type TelegramQueryResponse struct {
	Result telebot.Result
}

type TelegramCallbackResponse struct {
	Response *telebot.CallbackResponse
}
