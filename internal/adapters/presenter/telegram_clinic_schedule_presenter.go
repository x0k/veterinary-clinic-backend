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

const calendarInputValidationSchema = `{"type":"object","properties":{"selectedDates":{"type":"array","minItems":1}},"required":["selectedDates"]}`

type telegramClinicSchedulePresenter struct {
	calendarWebAppUrl           adapters.CalendarWebAppUrl
	calendarInputRequestOptions string
}

func newTelegramClinicSchedulePresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) telegramClinicSchedulePresenter {
	return telegramClinicSchedulePresenter{
		calendarWebAppUrl:           calendarWebAppUrl,
		calendarInputRequestOptions: fmt.Sprintf(`{"url":"%s"}`, string(calendarWebHandlerUrl)),
	}
}

func (p *telegramClinicSchedulePresenter) scheduleButtons(schedule entity.Schedule) []telebot.InlineButton {
	buttons := make([]telebot.InlineButton, 0, 3)
	if schedule.PrevDate != nil {
		buttons = append(buttons, *adapters.PreviousClinicScheduleBtn.With(schedule.PrevDate.Format(time.DateOnly)))
	}
	webAppParams := url.Values{}
	webAppParams.Add("r", p.calendarInputRequestOptions)
	webAppParams.Add("v", calendarInputValidationSchema)
	webAppParams.Add("w", fmt.Sprintf(
		`{"date":{"min":"%s"},"settings":{"selected":{"dates":["%s"]}}}`,
		time.Now().Format(time.DateOnly),
		schedule.Date.Format(time.DateOnly),
	))
	url := fmt.Sprintf("%s?%s", p.calendarWebAppUrl, webAppParams.Encode())
	buttons = append(buttons, telebot.InlineButton{
		Text: "📅",
		WebApp: &telebot.WebApp{
			URL: url,
		},
	})
	if schedule.NextDate != nil {
		buttons = append(buttons, *adapters.NextClinicScheduleBtn.With(schedule.NextDate.Format(time.DateOnly)))
	}
	return buttons
}

func (p *telegramClinicSchedulePresenter) schedule(schedule entity.Schedule) string {
	sb := strings.Builder{}
	sb.WriteString("График работы на ")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(
		schedule.Date.Format("02.01.2006")),
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
	return sb.String()
}

type TelegramClinicScheduleTextPresenter struct {
	telegramClinicSchedulePresenter
}

func NewTelegramClinicScheduleTextPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) *TelegramClinicScheduleTextPresenter {
	return &TelegramClinicScheduleTextPresenter{
		telegramClinicSchedulePresenter: newTelegramClinicSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *TelegramClinicScheduleTextPresenter) RenderSchedule(schedule entity.Schedule) (adapters.TelegramTextResponse, error) {
	return adapters.TelegramTextResponse{
		Text: p.schedule(schedule),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					p.scheduleButtons(schedule),
				},
			},
		},
	}, nil
}

type TelegramClinicScheduleQueryPresenter struct {
	telegramClinicSchedulePresenter
}

func NewTelegramClinicScheduleQueryPresenter(
	calendarWebAppUrl adapters.CalendarWebAppUrl,
	calendarWebHandlerUrl adapters.CalendarWebHandlerUrl,
) *TelegramClinicScheduleQueryPresenter {
	return &TelegramClinicScheduleQueryPresenter{
		telegramClinicSchedulePresenter: newTelegramClinicSchedulePresenter(
			calendarWebAppUrl,
			calendarWebHandlerUrl,
		),
	}
}

func (p *TelegramClinicScheduleQueryPresenter) RenderSchedule(schedule entity.Schedule) (adapters.TelegramQueryResponse, error) {
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
				ReplyMarkup: &telebot.ReplyMarkup{
					InlineKeyboard: [][]telebot.InlineButton{
						p.scheduleButtons(schedule),
					},
				},
			},
			Title: "График работы",
			Text:  p.schedule(schedule),
		},
	}, nil
}
