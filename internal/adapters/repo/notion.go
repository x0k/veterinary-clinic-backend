package repo

import (
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
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
	RecordAwaits    = "Ожидает"
	RecordInWork    = "В работе"
	RecordDone      = "Выполнено"
	RecordNotAppear = "Не пришел"
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

func Relations(properties notionapi.Properties, relationKey string) []notionapi.Relation {
	return properties[relationKey].(*notionapi.RelationProperty).Relation
}

func Service(page notionapi.Page) entity.Service {
	return entity.Service{
		Id:                entity.ServiceId(page.ID),
		Title:             Title(page.Properties, ServiceTitle),
		DurationInMinutes: entity.DurationInMinutes(Number(page.Properties, ServiceDurationInMinutes)),
		Description:       Text(page.Properties, ServiceDescription),
		CostDescription:   Text(page.Properties, ServiceCost),
	}
}

func UserIdFromRecord(properties notionapi.Properties, currentUserId *entity.UserId) *entity.UserId {
	if currentUserId == nil {
		return nil
	}
	uid := Text(properties, RecordUserId)
	if uid != string(*currentUserId) {
		return nil
	}
	return currentUserId
}

func ActualRecordStatus(properties notionapi.Properties) entity.RecordStatus {
	status := properties[RecordState].(*notionapi.SelectProperty).Select.Name
	if status == RecordInWork {
		return entity.RecordInWork
	}
	return entity.RecordAwaits
}

func DateTimePeriod(properties notionapi.Properties, key string) *entity.DateTimePeriod {
	date := Date(properties, key)
	if date == nil {
		return nil
	}
	start := date.Start
	if start == nil {
		return nil
	}
	end := date.End
	if end == nil {
		return nil
	}
	return &entity.DateTimePeriod{
		Start: entity.GoTimeToDateTime(time.Time(*start)),
		End:   entity.GoTimeToDateTime(time.Time(*end)),
	}
}

func DateTimePeriodFromRecord(properties notionapi.Properties) *entity.DateTimePeriod {
	return DateTimePeriod(properties, RecordDateTimePeriod)
}

func ActualRecord(page notionapi.Page, currentUserId *entity.UserId, service entity.Service) *entity.Record {
	dateTimePeriod := DateTimePeriodFromRecord(page.Properties)
	if dateTimePeriod == nil {
		return nil
	}
	return &entity.Record{
		Id:             entity.RecordId(page.ID),
		UserId:         UserIdFromRecord(page.Properties, currentUserId),
		Status:         ActualRecordStatus(page.Properties),
		DateTimePeriod: *dateTimePeriod,
		Service:        service,
	}
}

func PrivateActualRecord(page notionapi.Page, service entity.Service) (entity.Record, error) {
	dateTimePeriod := DateTimePeriodFromRecord(page.Properties)
	if dateTimePeriod == nil {
		return entity.Record{}, ErrFailedToCreateRecord
	}
	uid := entity.UserId(Text(page.Properties, RecordUserId))
	return entity.Record{
		Id:             entity.RecordId(page.ID),
		UserId:         &uid,
		Status:         ActualRecordStatus(page.Properties),
		DateTimePeriod: *dateTimePeriod,
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

var breakTimePeriod = entity.TimePeriod{
	Start: entity.Time{
		Hours:   0,
		Minutes: 0,
	},
	End: entity.Time{
		Hours:   23,
		Minutes: 59,
	},
}

func WorkBreak(page notionapi.Page) *entity.WorkBreak {
	period := DateTimePeriod(page.Properties, BreakPeriod)
	if period == nil {
		return nil
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
	return &entity.WorkBreak{
		Id:              entity.WorkBreakId(page.ID),
		Title:           Title(page.Properties, BreakTitle),
		MatchExpression: sb.String(),
		Period:          breakTimePeriod,
	}
}
