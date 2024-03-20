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

func (s *NotionServicesRepo) Service(ctx context.Context, serviceId entity.ServiceId) (entity.Service, error) {
	r, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return entity.Service{}, err
	}
	if r == nil {
		return entity.Service{}, ErrNotFound
	}
	return Service(*r), nil
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
