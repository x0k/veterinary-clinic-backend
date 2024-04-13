package repo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const (
	ServiceTitle             = "Наименование"
	ServiceDurationInMinutes = "Продолжительность в минутах"
	ServiceDescription       = "Описание"
	ServiceCost              = "Стоимость"
)

const (
	RecordTitle          = "ФИО"
	RecordService        = "Услуга"
	RecordPhoneNumber    = "Телефон"
	RecordEmail          = "Почта"
	RecordDateTimePeriod = "Время записи"
	RecordState          = "Статус"
	RecordUserId         = "identity"
)

const (
	BreakTitle  = "Наименование"
	BreakPeriod = "Период"
)

const (
	RecordAwaits            = "Ожидает"
	RecordDone              = "Выполнено"
	RecordNotAppear         = "Не пришел"
	RecordDoneArchived      = "Архив выполнено"
	RecordNotAppearArchived = "Архив не пришел"
)

func RichTextValue(richText []notionapi.RichText) string {
	strs := make([]string, 0, len(richText))
	for _, r := range richText {
		if r.Type != notionapi.ObjectTypeText {
			continue
		}
		strs = append(strs, r.Text.Content)
	}
	return strings.Join(strs, "")
}

func Title(properties notionapi.Properties, titleKey string) string {
	return RichTextValue(properties[titleKey].(*notionapi.TitleProperty).Title)
}

func Number(properties notionapi.Properties, numberKey string) float64 {
	return properties[numberKey].(*notionapi.NumberProperty).Number
}

func Text(properties notionapi.Properties, stringKey string) string {
	return RichTextValue(properties[stringKey].(*notionapi.RichTextProperty).RichText)
}

func Date(properties notionapi.Properties, dateKey string) *notionapi.DateObject {
	return properties[dateKey].(*notionapi.DateProperty).Date
}

func Phone(properties notionapi.Properties, phoneKey string) string {
	return properties[phoneKey].(*notionapi.PhoneNumberProperty).PhoneNumber
}

func Email(properties notionapi.Properties, emailKey string) string {
	return properties[emailKey].(*notionapi.EmailProperty).Email
}

func Relations(properties notionapi.Properties, relationKey string) []notionapi.Relation {
	return properties[relationKey].(*notionapi.RelationProperty).Relation
}

func Service(page notionapi.Page) shared.Service {
	return shared.Service{
		Id:                shared.ServiceId(page.ID),
		Title:             Title(page.Properties, ServiceTitle),
		DurationInMinutes: shared.DurationInMinutes(Number(page.Properties, ServiceDurationInMinutes)),
		Description:       Text(page.Properties, ServiceDescription),
		CostDescription:   Text(page.Properties, ServiceCost),
	}
}

func RecordStatus(properties notionapi.Properties) (shared.RecordStatus, error) {
	switch properties[RecordState].(*notionapi.SelectProperty).Select.Name {
	case RecordAwaits:
		return shared.RecordAwaits, nil
	case RecordDone:
		return shared.RecordDone, nil
	case RecordNotAppear:
		return shared.RecordNotAppear, nil
	case RecordDoneArchived:
		return shared.RecordDoneArchived, nil
	case RecordNotAppearArchived:
		return shared.RecordNotAppearArchived, nil
	default:
		return shared.RecordStatus(""), shared.ErrInvalidRecordStatus
	}
}

func DateTimePeriod(properties notionapi.Properties, key string) (shared.DateTimePeriod, error) {
	date := Date(properties, key)
	if date == nil || date.Start == nil || date.End == nil {
		return shared.DateTimePeriod{}, fmt.Errorf("%s: %w", key, shared.ErrInvalidDate)
	}
	return shared.DateTimePeriod{
		Start: shared.GoTimeToDateTime(time.Time(*date.Start)),
		End:   shared.GoTimeToDateTime(time.Time(*date.End)),
	}, nil
}

func User(properties notionapi.Properties) shared.User {
	return shared.User{
		Id:          shared.UserId(Text(properties, RecordUserId)),
		Name:        Title(properties, RecordTitle),
		PhoneNumber: Phone(properties, RecordPhoneNumber),
		Email:       Email(properties, RecordEmail),
	}
}

func Record(page notionapi.Page, service shared.Service) (shared.Record, error) {
	dateTimePeriod, err := DateTimePeriod(page.Properties, RecordDateTimePeriod)
	if err != nil {
		return shared.Record{}, err
	}
	status, err := RecordStatus(page.Properties)
	if err != nil {
		return shared.Record{}, err
	}
	return shared.Record{
		Id:             shared.RecordId(page.ID),
		User:           User(page.Properties),
		Status:         status,
		DateTimePeriod: dateTimePeriod,
		Service:        service,
	}, nil
}

func RichText(value string) []notionapi.RichText {
	return []notionapi.RichText{
		{
			Type: notionapi.ObjectTypeText,
			Text: &notionapi.Text{Content: value},
		},
	}
}

func WorkBreak(page notionapi.Page) (shared.WorkBreak, error) {
	period, err := DateTimePeriod(page.Properties, BreakPeriod)
	if err != nil {
		return shared.WorkBreak{}, err
	}
	dt := time.Date(
		period.Start.Year,
		time.Month(period.Start.Month),
		period.Start.Day,
		0, 0, 0, 0, time.Local)
	sb := strings.Builder{}
	sb.WriteString("^\\d (")
	for dt.Year() < period.End.Year ||
		dt.Month() < time.Month(period.End.Month) ||
		dt.Day() < period.End.Day {
		sb.WriteString(dt.Format(time.DateOnly))
		sb.WriteByte('|')
		dt = dt.AddDate(0, 0, 1)
	}
	sb.WriteString(dt.Format(time.DateOnly))
	sb.WriteByte(')')
	return shared.WorkBreak{
		Id:              shared.WorkBreakId(page.ID),
		Title:           Title(page.Properties, BreakTitle),
		MatchExpression: sb.String(),
		Period: shared.TimePeriod{
			Start: period.Start.Time,
			End:   period.End.Time,
		},
	}, nil
}
