package shared

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type TelegramResponse struct {
	Text    string
	Options *telebot.SendOptions
}

var escapeRegExp = regexp.MustCompile(`([.!-])`)

func EscapeTelegramMarkdownString(text string) string {
	return escapeRegExp.ReplaceAllString(text, "\\$1")
}
