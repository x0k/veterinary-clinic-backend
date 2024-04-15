package telegram_adapters

import "gopkg.in/telebot.v3"

type TextMessages struct {
	recipient telebot.Recipient
	messages  []SendableText
}

func NewTextMessages(
	recipient telebot.Recipient,
	messages ...SendableText,
) TextMessages {
	return TextMessages{
		recipient: recipient,
		messages:  messages,
	}
}

func (m TextMessages) Send(bot *telebot.Bot) error {
	for _, msg := range m.messages {
		if _, err := bot.Send(
			m.recipient,
			msg.Text,
			msg.Options,
		); err != nil {
			return err
		}
	}
	return nil
}
