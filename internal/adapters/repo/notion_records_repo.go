package repo

import (
	"context"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type NotionRecordsRepo struct {
	recordsDatabaseId notionapi.DatabaseID
	client            *notionapi.Client
}

func NewNotionRecords(
	client *notionapi.Client,
	recordsDatabaseId notionapi.DatabaseID,
) *NotionRecordsRepo {
	return &NotionRecordsRepo{
		client:            client,
		recordsDatabaseId: recordsDatabaseId,
	}
}

func (s *NotionRecordsRepo) Create(
	ctx context.Context,
	user entity.User,
	service entity.Service,
	appointmentDateTime time.Time,
) (entity.Record, error) {
	start := notionapi.Date(appointmentDateTime)
	end := notionapi.Date(
		appointmentDateTime.Add(time.Duration(service.DurationInMinutes) * time.Minute),
	)
	properties := notionapi.Properties{
		RecordTitle: notionapi.TitleProperty{
			Type:  notionapi.PropertyTypeTitle,
			Title: RichText(user.Name),
		},
		RecordService: notionapi.RelationProperty{
			Type: notionapi.PropertyTypeRelation,
			Relation: []notionapi.Relation{
				{
					ID: notionapi.PageID(service.Id),
				},
			},
		},
		RecordEmail: notionapi.EmailProperty{
			Type:  notionapi.PropertyTypeEmail,
			Email: user.Email,
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
			Type:     notionapi.PropertyTypeRichText,
			RichText: RichText(string(user.Id)),
		},
	}
	if user.PhoneNumber != "" {
		properties[RecordPhoneNumber] = notionapi.PhoneNumberProperty{
			Type:        notionapi.PropertyTypePhoneNumber,
			PhoneNumber: user.PhoneNumber,
		}
	}
	res, err := s.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: s.recordsDatabaseId,
		},
		Properties: properties,
	})
	if err != nil {
		return entity.Record{}, err
	}
	if res == nil {
		return entity.Record{}, ErrFailedToCreateRecord
	}
	if rec := ActualRecord(*res, &user.Id, service); rec != nil {
		return *rec, nil
	}
	return entity.Record{}, ErrFailedToCreateRecord
}

func (s *NotionRecordsRepo) Remove(ctx context.Context, recordId entity.RecordId) error {
	_, err := s.client.Page.Update(ctx, notionapi.PageID(recordId), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{},
		Archived:   true,
	})
	return err
}

func (s *NotionRecordsRepo) RecordByUserId(ctx context.Context, userId entity.UserId) (entity.Record, error) {
	res, err := s.recordDbRespByUserId(ctx, userId)
	if err != nil {
		return entity.Record{}, err
	}
	if len(res.Results) == 0 {
		return entity.Record{}, usecase.ErrNotFound
	}
	page := res.Results[0]
	relations := Relations(page.Properties, RecordService)
	if len(relations) == 0 {
		return entity.Record{}, adapters.ErrInvalidRecord
	}
	service, err := s.client.Page.Get(ctx, notionapi.PageID(relations[0].ID))
	if err != nil {
		return entity.Record{}, err
	}
	if service == nil {
		return entity.Record{}, usecase.ErrNotFound
	}
	return *ActualRecord(page, &userId, Service(*service)), nil
}

func (s *NotionRecordsRepo) recordDbRespByUserId(ctx context.Context, userId entity.UserId) (*notionapi.DatabaseQueryResponse, error) {
	return s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
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
			notionapi.PropertyFilter{
				Property: RecordUserId,
				RichText: &notionapi.TextFilterCondition{
					Equals: string(userId),
				},
			},
		},
	})
}
