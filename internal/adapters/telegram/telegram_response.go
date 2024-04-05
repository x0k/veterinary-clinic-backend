package adapters_telegram

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type TextResponse struct {
	Text    string
	Options *telebot.SendOptions
}

var escapeRegExp = regexp.MustCompile(`([.!-])`)

func EscapeMarkdownString(text string) string {
	return escapeRegExp.ReplaceAllString(text, "\\$1")
}

type QueryResponse struct {
	Result telebot.Result
}

type CallbackResponse struct {
	Response *telebot.CallbackResponse
}
