package repo

import (
	"context"
	"errors"
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

var ErrFailedToCreateRecord = errors.New("failed to create record")

type Notion struct {
	servicesDatabaseId                notionapi.DatabaseID
	recordsDatabaseId                 notionapi.DatabaseID
	client                            *notionapi.Client
	actualRecordsDatabaseQueryRequest *notionapi.DatabaseQueryRequest
}

func (s *Notion) richTextValue(richText []notionapi.RichText) string {
	strs := make([]string, 0, len(richText))
	for _, r := range richText {
		if r.Type != notionapi.ObjectTypeText {
			continue
		}
		strs = append(strs, r.Text.Content)
	}
	return strings.Join(strs, "")
}

func (s *Notion) title(properties notionapi.Properties, titleKey string) string {
	return s.richTextValue(properties[titleKey].(*notionapi.TitleProperty).Title)
}

func (s *Notion) number(properties notionapi.Properties, numberKey string) float64 {
	return properties[numberKey].(*notionapi.NumberProperty).Number
}

func (s *Notion) text(properties notionapi.Properties, stringKey string) string {
	return s.richTextValue(properties[stringKey].(*notionapi.RichTextProperty).RichText)
}

func (s *Notion) date(properties notionapi.Properties, dateKey string) *notionapi.DateObject {
	return properties[dateKey].(*notionapi.DateProperty).Date
}

func (s *Notion) service(page notionapi.Page) entity.Service {
	return entity.Service{
		Id:                entity.ServiceId(page.ID),
		Title:             s.title(page.Properties, ServiceTitle),
		DurationInMinutes: int(s.number(page.Properties, ServiceDurationInMinutes)),
		Description:       s.text(page.Properties, ServiceDescription),
		CostDescription:   s.text(page.Properties, ServiceCost),
	}
}

func (s *Notion) recordUserId(properties notionapi.Properties, currentUserId *entity.UserId) *entity.UserId {
	if currentUserId == nil {
		return nil
	}
	uid := s.text(properties, RecordUserId)
	if uid != string(*currentUserId) {
		return nil
	}
	return currentUserId
}

func (s *Notion) actualRecordStatus(properties notionapi.Properties) entity.RecordStatus {
	status := properties[RecordState].(*notionapi.SelectProperty).Select.Name
	if status == ClinicRecordInWork {
		return entity.RecordInWork
	}
	return entity.RecordAwaits
}

func (s *Notion) actualRecord(page notionapi.Page, currentUserId *entity.UserId) *entity.Record {
	startDateTime := s.date(page.Properties, RecordDateTimePeriod).Start
	if startDateTime == nil {
		return nil
	}
	endDateTime := s.date(page.Properties, RecordDateTimePeriod).End
	if endDateTime == nil {
		return nil
	}
	return &entity.Record{
		Id:     entity.RecordId(page.ID),
		UserId: s.recordUserId(page.Properties, currentUserId),
		Status: s.actualRecordStatus(page.Properties),
		DateTimePeriod: entity.DateTimePeriod{
			Start: entity.GoTimeToDateTime(time.Time(*startDateTime)),
			End:   entity.GoTimeToDateTime(time.Time(*endDateTime)),
		},
	}
}

func (s *Notion) makeRichText(value string) []notionapi.RichText {
	return []notionapi.RichText{
		{
			Type: notionapi.ObjectTypeText,
			Text: &notionapi.Text{Content: value},
		},
	}
}

func New(
	client *notionapi.Client,
	servicesDatabaseId notionapi.DatabaseID,
	recordsDatabaseId notionapi.DatabaseID,
) *Notion {
	return &Notion{
		client:             client,
		servicesDatabaseId: servicesDatabaseId,
		recordsDatabaseId:  recordsDatabaseId,
		actualRecordsDatabaseQueryRequest: &notionapi.DatabaseQueryRequest{
			Filter: notionapi.AndCompoundFilter{
				notionapi.PropertyFilter{
					Property: RecordDateTimePeriod,
					Date: &notionapi.DateFilterCondition{
						IsNotEmpty: true,
					},
				},
				notionapi.OrCompoundFilter{
					notionapi.PropertyFilter{
						Property: RecordState,
						Select: &notionapi.SelectFilterCondition{
							Equals: ClinicRecordInWork,
						},
					},
					notionapi.PropertyFilter{
						Property: RecordState,
						Select: &notionapi.SelectFilterCondition{
							Equals: ClinicRecordAwaits,
						},
					},
				},
			},
			Sorts: []notionapi.SortObject{
				{
					Property:  RecordDateTimePeriod,
					Direction: notionapi.SortOrderASC,
				},
			},
		},
	}
}

func (s *Notion) Services(ctx context.Context) ([]entity.Service, error) {
	r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
	if err != nil {
		return nil, err
	}
	services := make([]entity.Service, 0, len(r.Results))
	for _, result := range r.Results {
		services = append(services, s.service(result))
	}
	return services, nil
}

func (s *Notion) FetchActualRecords(ctx context.Context, currentUserId *entity.UserId) ([]entity.Record, error) {
	r, err := s.client.Database.Query(ctx, s.recordsDatabaseId, s.actualRecordsDatabaseQueryRequest)
	if err != nil {
		return nil, err
	}
	records := make([]entity.Record, 0, len(r.Results))
	for _, result := range r.Results {
		if rec := s.actualRecord(result, currentUserId); rec != nil {
			records = append(records, *rec)
		}
	}
	return records, nil
}

func (s *Notion) Records(ctx context.Context) ([]entity.Record, error) {
	return s.FetchActualRecords(ctx, nil)
}

func (s *Notion) CreateRecord(
	ctx context.Context,
	userId entity.UserId,
	serviceId entity.ServiceId,
	userName string,
	userEmail string,
	userPhoneNumber string,
	utcDateTimePeriod entity.DateTimePeriod,
) (entity.Record, error) {
	start := notionapi.Date(entity.DateTimeToGoTime(utcDateTimePeriod.Start))
	end := notionapi.Date(entity.DateTimeToGoTime(utcDateTimePeriod.End))
	res, err := s.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: s.recordsDatabaseId,
		},
		Properties: notionapi.Properties{
			RecordTitle: notionapi.TitleProperty{
				Type:  notionapi.PropertyTypeTitle,
				Title: s.makeRichText(userName),
			},
			RecordService: notionapi.RelationProperty{
				Type: notionapi.PropertyTypeRelation,
				Relation: []notionapi.Relation{
					{
						ID: notionapi.PageID(serviceId),
					},
				},
			},
			RecordPhoneNumber: notionapi.PhoneNumberProperty{
				Type:        notionapi.PropertyTypePhoneNumber,
				PhoneNumber: userPhoneNumber,
			},
			RecordEmail: notionapi.EmailProperty{
				Type:  notionapi.PropertyTypeEmail,
				Email: userEmail,
			},
			RecordDateTimePeriod: notionapi.DateProperty{
				Type: notionapi.PropertyTypeDate,
				Date: &notionapi.DateObject{
					Start: &start,
					End:   &end,
				},
			},
			RecordState: notionapi.SelectProperty{
				Type: notionapi.PropertyTypeSelect,
				Select: notionapi.Option{
					Name: ClinicRecordAwaits,
				},
			},
			RecordUserId: notionapi.RichTextProperty{
				Type:     notionapi.PropertyTypeText,
				RichText: s.makeRichText(string(userId)),
			},
		},
	})
	if err != nil {
		return entity.Record{}, err
	}
	if res == nil {
		return entity.Record{}, ErrFailedToCreateRecord
	}
	if rec := s.actualRecord(*res, &userId); rec != nil {
		return *rec, nil
	}
	return entity.Record{}, ErrFailedToCreateRecord
}

func (s *Notion) RemoveRecord(ctx context.Context, recordId entity.RecordId) error {
	_, err := s.client.Page.Update(ctx, notionapi.PageID(recordId), &notionapi.PageUpdateRequest{
		Archived: true,
	})
	return err
}
