package repo

import (
	"context"
	"errors"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrFailedToCreateRecord = errors.New("failed to create record")

type NotionServicesRepo struct {
	servicesDatabaseId                notionapi.DatabaseID
	recordsDatabaseId                 notionapi.DatabaseID
	client                            *notionapi.Client
	actualRecordsDatabaseQueryRequest *notionapi.DatabaseQueryRequest
}

func NewNotionServices(
	client *notionapi.Client,
	servicesDatabaseId notionapi.DatabaseID,
	recordsDatabaseId notionapi.DatabaseID,
) *NotionServicesRepo {
	return &NotionServicesRepo{
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
							Equals: RecordInWork,
						},
					},
					notionapi.PropertyFilter{
						Property: RecordState,
						Select: &notionapi.SelectFilterCondition{
							Equals: RecordAwaits,
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

func (s *NotionServicesRepo) Services(ctx context.Context) ([]entity.Service, error) {
	r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
	if err != nil {
		return nil, err
	}
	services := make([]entity.Service, 0, len(r.Results))
	for _, result := range r.Results {
		services = append(services, Service(result))
	}
	return services, nil
}

func (s *NotionServicesRepo) FetchActualRecords(ctx context.Context, currentUserId *entity.UserId) ([]entity.Record, error) {
	r, err := s.client.Database.Query(ctx, s.recordsDatabaseId, s.actualRecordsDatabaseQueryRequest)
	if err != nil {
		return nil, err
	}
	records := make([]entity.Record, 0, len(r.Results))
	for _, result := range r.Results {
		if rec := ActualRecord(result, currentUserId); rec != nil {
			records = append(records, *rec)
		}
	}
	return records, nil
}

func (s *NotionServicesRepo) CreateRecord(
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
				Title: RichText(userName),
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
					Name: RecordAwaits,
				},
			},
			RecordUserId: notionapi.RichTextProperty{
				Type:     notionapi.PropertyTypeText,
				RichText: RichText(string(userId)),
			},
		},
	})
	if err != nil {
		return entity.Record{}, err
	}
	if res == nil {
		return entity.Record{}, ErrFailedToCreateRecord
	}
	if rec := ActualRecord(*res, &userId); rec != nil {
		return *rec, nil
	}
	return entity.Record{}, ErrFailedToCreateRecord
}

func (s *NotionServicesRepo) RemoveRecord(ctx context.Context, recordId entity.RecordId) error {
	_, err := s.client.Page.Update(ctx, notionapi.PageID(recordId), &notionapi.PageUpdateRequest{
		Archived: true,
	})
	return err
}
