package presenter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramDialogPresenter struct {
	calendarWebAppUrl           string
	calendarInputRequestOptions string
}

const calendarInputValidationSchema = `{"type":"object","properties":{"selectedDates":{"type":"array","minItems":1}},"required":["selectedDates"]}`

func NewTelegramDialog(calendarWebAppUrl string, calendarWebHandlerUrl string) *TelegramDialogPresenter {
	return &TelegramDialogPresenter{
		calendarWebAppUrl:           calendarWebAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url": "%s"}`, calendarWebHandlerUrl),
	}
}

func (p *TelegramDialogPresenter) RenderGreeting() (adapters.TelegramResponse, error) {
	return adapters.TelegramTextResponse{
		Text: adapters.EscapeTelegramMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderDatePicker(t time.Time) (adapters.TelegramResponse, error) {
	webAppParams := url.Values{}
	webAppParams.Add("r", p.calendarInputRequestOptions)
	webAppParams.Add("v", calendarInputValidationSchema)
	date := t.Format(time.DateOnly)
	webAppParams.Add("w", fmt.Sprintf(`{"date":{"min":"%s"},"settings":{"selected":{"dates":["%s"]}}}`, date, date))
	url := fmt.Sprintf("%s?%s", p.calendarWebAppUrl, webAppParams.Encode())
	return adapters.TelegramTextResponse{
		Text: "Выберите дату",
		Options: &telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{{
					{
						Text: "Открыть календарь",
						WebApp: &telebot.WebApp{
							URL: url,
						},
					},
				}},
			},
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderSchedule(schedule entity.Schedule) (adapters.TelegramResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("График работы на ")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(
		entity.DateToGoTime(schedule.Date).Format("02.01.2006")),
	)
	sb.WriteString(":\n\n")
	for _, period := range schedule.Periods {
		sb.WriteByte('*')
		sb.WriteString(period.Start.String())
		sb.WriteString(" \\- ")
		sb.WriteString(period.End.String())
		sb.WriteString("*\n")
		sb.WriteString(adapters.EscapeTelegramMarkdownString(period.Title))
		sb.WriteString("\n\n")
	}
	if len(schedule.Periods) == 0 {
		sb.WriteString("Нет записей\n\n")
	}
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
			},
			Title: "График работы",
			Text:  sb.String(),
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderError(err error) (adapters.TelegramResponse, error) {
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &err),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
			},
			Title: "Ошибка",
			Text:  adapters.EscapeTelegramMarkdownString(err.Error()),
		},
	}, nil
}
