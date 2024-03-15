package notion_repo

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
	ClinicRecordAwaits    = "Ожидает"
	ClinicRecordInWork    = "В работе"
	ClinicRecordDone      = "Выполнено"
	ClinicRecordNotAppear = "Не пришел"
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
		DurationInMinutes: int(Number(page.Properties, ServiceDurationInMinutes)),
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
	if status == ClinicRecordInWork {
		return entity.RecordInWork
	}
	return entity.RecordAwaits
}

func ActualRecord(page notionapi.Page, currentUserId *entity.UserId) *entity.Record {
	startDateTime := Date(page.Properties, RecordDateTimePeriod).Start
	if startDateTime == nil {
		return nil
	}
	endDateTime := Date(page.Properties, RecordDateTimePeriod).End
	if endDateTime == nil {
		return nil
	}
	return &entity.Record{
		Id:     entity.RecordId(page.ID),
		UserId: UserIdFromRecord(page.Properties, currentUserId),
		Status: ActualRecordStatus(page.Properties),
		DateTimePeriod: entity.DateTimePeriod{
			Start: entity.GoTimeToDateTime(time.Time(*startDateTime)),
			End:   entity.GoTimeToDateTime(time.Time(*endDateTime)),
		},
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
