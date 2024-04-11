package adapters_telegram

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

type textResponse struct {
	Text    string
	Options *telebot.SendOptions
}

type TextResponses = []textResponse

func Send(c telebot.Context, responses TextResponses) error {
	for _, response := range responses {
		if err := c.Send(response.Text, response.Options); err != nil {
			return err
		}
	}
	return nil
}

func Edit(c telebot.Context, responses TextResponses) error {
	if len(responses) == 0 {
		return nil
	}
	c.Edit(responses[0].Text, responses[0].Options)
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
