package repo

import (
	"context"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
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
	if rec := ActualRecord(*res, &user.Id); rec != nil {
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
