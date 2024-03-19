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

func DateTimePeriodFromRecord(properties notionapi.Properties) *entity.DateTimePeriod {
	date := Date(properties, RecordDateTimePeriod)
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

func ActualRecord(page notionapi.Page, currentUserId *entity.UserId) *entity.Record {
	dateTimePeriod := DateTimePeriodFromRecord(page.Properties)
	if dateTimePeriod == nil {
		return nil
	}
	return &entity.Record{
		Id:             entity.RecordId(page.ID),
		UserId:         UserIdFromRecord(page.Properties, currentUserId),
		Status:         ActualRecordStatus(page.Properties),
		DateTimePeriod: *dateTimePeriod,
	}
}

func RichText(value string) []notionapi.RichText {
	return []notionapi.RichText{
		{
			Type: notionapi.ObjectTypeText,
			Text: &notionapi.Text{Content: value},
		},
	}
}
