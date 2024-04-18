package telegram_adapters

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type Message interface {
	Send(bot *telebot.Bot) error
}

type Response interface {
	Send(c telebot.Context) error
}

type SendableText struct {
	Text    string
	Options *telebot.SendOptions
}

func NewSendableText(
	text string,
	options ...*telebot.SendOptions,
) SendableText {
	if len(options) == 0 {
		return SendableText{
			Text:    text,
			Options: &telebot.SendOptions{},
		}
	}
	return SendableText{
		Text:    text,
		Options: options[0],
	}
}

type TextResponses []SendableText

func (rs TextResponses) Send(c telebot.Context) error {
	for _, response := range rs {
		if err := c.Send(response.Text, response.Options); err != nil {
			return err
		}
	}
	return nil
}

func (responses TextResponses) Edit(c telebot.Context) error {
	if len(responses) == 0 {
		return nil
	}
	if err := c.Edit(responses[0].Text, responses[0].Options); err != nil {
		return err
	}
	for i := 1; i < len(responses); i++ {
		if err := c.Reply(responses[i].Text, responses[i].Options); err != nil {
			return err
		}
	}
	return nil
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
