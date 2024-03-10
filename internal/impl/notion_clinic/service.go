package notion_clinic

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/models"
)

var ErrFailedToCreateRecord = errors.New("failed to create record")

type Service struct {
	servicesDatabaseId                notionapi.DatabaseID
	recordsDatabaseId                 notionapi.DatabaseID
	client                            *notionapi.Client
	actualRecordsDatabaseQueryRequest *notionapi.DatabaseQueryRequest
}

func (s *Service) richTextValue(richText []notionapi.RichText) string {
	strs := make([]string, 0, len(richText))
	for _, r := range richText {
		if r.Type != notionapi.ObjectTypeText {
			continue
		}
		strs = append(strs, r.Text.Content)
	}
	return strings.Join(strs, "")
}

func (s *Service) title(properties notionapi.Properties, titleKey string) string {
	return s.richTextValue(properties[titleKey].(notionapi.TitleProperty).Title)
}

func (s *Service) number(properties notionapi.Properties, numberKey string) float64 {
	return properties[numberKey].(notionapi.NumberProperty).Number
}

func (s *Service) text(properties notionapi.Properties, stringKey string) string {
	return s.richTextValue(properties[stringKey].(notionapi.TextProperty).Text)
}

func (s *Service) date(properties notionapi.Properties, dateKey string) *notionapi.DateObject {
	return properties[dateKey].(notionapi.DateProperty).Date
}

func (s *Service) service(page notionapi.Page) models.Service {
	return models.Service{
		Id:                models.ServiceId(page.ID),
		Title:             s.title(page.Properties, ServiceTitle),
		DurationInMinutes: int(s.number(page.Properties, ServiceDurationInMinutes)),
		Description:       s.text(page.Properties, ServiceDescription),
		CostDescription:   s.text(page.Properties, ServiceCost),
	}
}

func (s *Service) recordUserId(properties notionapi.Properties, currentUserId *models.UserId) *models.UserId {
	if currentUserId == nil {
		return nil
	}
	uid := s.text(properties, RecordUserId)
	if uid != string(*currentUserId) {
		return nil
	}
	return currentUserId
}

func (s *Service) actualRecordStatus(properties notionapi.Properties) models.RecordStatus {
	status := properties[RecordState].(notionapi.SelectProperty).Select.Name
	if status == ClinicRecordInWork {
		return models.RecordInWork
	}
	return models.RecordAwaits
}

func (s *Service) actualRecord(page notionapi.Page, currentUserId *models.UserId) *models.Record {
	startDateTime := s.date(page.Properties, RecordDateTimePeriod).Start
	if startDateTime == nil {
		return nil
	}
	endDateTime := s.date(page.Properties, RecordDateTimePeriod).End
	if endDateTime == nil {
		return nil
	}
	return &models.Record{
		Id:     models.RecordId(page.ID),
		UserId: s.recordUserId(page.Properties, currentUserId),
		Status: s.actualRecordStatus(page.Properties),
		DateTimePeriod: models.DateTimePeriod{
			Start: models.GoTimeToDateTime(time.Time(*startDateTime)),
			End:   models.GoTimeToDateTime(time.Time(*endDateTime)),
		},
	}
}

func (s *Service) makeRichText(value string) []notionapi.RichText {
	return []notionapi.RichText{
		{
			Type: notionapi.ObjectTypeText,
			Text: &notionapi.Text{Content: value},
		},
	}
}

func NewService(
	client *notionapi.Client,
	servicesDatabaseId notionapi.DatabaseID,
	recordsDatabaseId notionapi.DatabaseID,
) *Service {
	return &Service{
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

func (s *Service) FetchServices(ctx context.Context) ([]models.Service, error) {
	r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
	if err != nil {
		return nil, err
	}
	services := make([]models.Service, 0, len(r.Results))
	for _, result := range r.Results {
		services = append(services, s.service(result))
	}
	return services, nil
}

func (s *Service) FetchActualRecords(ctx context.Context, currentUserId *models.UserId) ([]models.Record, error) {
	r, err := s.client.Database.Query(ctx, s.recordsDatabaseId, s.actualRecordsDatabaseQueryRequest)
	if err != nil {
		return nil, err
	}
	records := make([]models.Record, 0, len(r.Results))
	for _, result := range r.Results {
		if rec := s.actualRecord(result, currentUserId); rec != nil {
			records = append(records, *rec)
		}
	}
	return records, nil
}

func (s *Service) CreateRecord(
	ctx context.Context,
	userId models.UserId,
	serviceId models.ServiceId,
	userName string,
	userEmail string,
	userPhoneNumber string,
	utcDateTimePeriod models.DateTimePeriod,
) (models.Record, error) {
	start := notionapi.Date(models.DateTimeToGoTime(utcDateTimePeriod.Start))
	end := notionapi.Date(models.DateTimeToGoTime(utcDateTimePeriod.End))
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
			RecordUserId: notionapi.TextProperty{
				Type: notionapi.PropertyTypeText,
				Text: s.makeRichText(string(userId)),
			},
		},
	})
	if err != nil {
		return models.Record{}, err
	}
	if res == nil {
		return models.Record{}, ErrFailedToCreateRecord
	}
	if rec := s.actualRecord(*res, &userId); rec != nil {
		return *rec, nil
	}
	return models.Record{}, ErrFailedToCreateRecord
}

func (s *Service) RemoveRecord(ctx context.Context, recordId models.RecordId) error {
	_, err := s.client.Page.Update(ctx, notionapi.PageID(recordId), &notionapi.PageUpdateRequest{
		Archived: true,
	})
	return err
}
