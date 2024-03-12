package shared

import "gopkg.in/telebot.v3"

type TelegramResponse struct {
	Text    string
	Options *telebot.SendOptions
}
