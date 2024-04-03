package appointment_notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

type AppointmentRepository struct {
	client            *notionapi.Client
	recordsDatabaseId notionapi.DatabaseID
}

func NewAppointment(
	client *notionapi.Client,
	recordsDatabaseId notionapi.DatabaseID,
) *AppointmentRepository {
	return &AppointmentRepository{
		client:            client,
		recordsDatabaseId: recordsDatabaseId,
	}
}

func (r *AppointmentRepository) IsAppointmentPeriodBusy(ctx context.Context, period entity.DateTimePeriod) (bool, error) {
	const op = "appointment_notion.AppointmentRepository.IsAppointmentPeriodBusy"
	after := notionapi.Date(entity.DateTimeToGoTime(period.Start))
	before := notionapi.Date(entity.DateTimeToGoTime(period.End))
	// TODO: Test notion api to query appointments with overlapping periods
	res, err := r.client.Database.Query(ctx, r.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: RecordDateTimePeriod,
				Date: &notionapi.DateFilterCondition{
					After: &after,
				},
			},
			notionapi.PropertyFilter{
				Property: RecordDateTimePeriod,
				Date: &notionapi.DateFilterCondition{
					Before: &before,
				},
			},
			notionapi.PropertyFilter{
				Property: RecordState,
				Select: &notionapi.SelectFilterCondition{
					Equals: RecordAwaits,
				},
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return len(res.Results) > 0, nil
}

func (r *AppointmentRepository) SaveAppointment(ctx context.Context, app *appointment.AppointmentAggregate) error {
	const op = "appointment_notion.AppointmentRepository.SaveAppointment"
	period := app.DateTimePeriod()
	start := notionapi.Date(entity.DateTimeToGoTime(period.Start))
	end := notionapi.Date(entity.DateTimeToGoTime(period.End))
	status, err := RecordStatusToNotion(app.State(), app.IsArchived())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	title, err := app.Title()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	properties := notionapi.Properties{
		RecordTitle: notionapi.TitleProperty{
			Type:  notionapi.PropertyTypeTitle,
			Title: notion.ToRichText(title),
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
				Name: status,
			},
		},
		RecordCustomer: notionapi.RelationProperty{
			Type: notionapi.PropertyTypeRelation,
			Relation: []notionapi.Relation{
				{
					ID: notionapi.PageID(app.CustomerId().String()),
				},
			},
		},
		RecordCreatedAt: notionapi.CreatedTimeProperty{
			Type:        notionapi.PropertyTypeCreatedTime,
			CreatedTime: app.CreatedAt(),
		},
		RecordService: notionapi.RelationProperty{
			Type: notionapi.PropertyTypeRelation,
			Relation: []notionapi.Relation{
				{
					ID: notionapi.PageID(app.ServiceId().String()),
				},
			},
		},
	}
	res, err := r.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: r.recordsDatabaseId,
		},
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	app.SetId(appointment.NewRecordId(string(res.ID)))
	return nil
}
