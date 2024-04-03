package appointment_notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

type RecordRepository struct {
	client                 *notionapi.Client
	appointmentsDatabaseId notionapi.DatabaseID
}

func NewAppointment(
	client *notionapi.Client,
	appointmentsDatabaseId notionapi.DatabaseID,
) *RecordRepository {
	return &RecordRepository{
		client:                 client,
		appointmentsDatabaseId: appointmentsDatabaseId,
	}
}
func (r *RecordRepository) IsAppointmentPeriodBusy(ctx context.Context, period entity.DateTimePeriod) (bool, error) {
	after := notionapi.Date(entity.DateTimeToGoTime(period.Start))
	before := notionapi.Date(entity.DateTimeToGoTime(period.End))
	res, err := r.client.Database.Query(ctx, r.appointmentsDatabaseId, &notionapi.DatabaseQueryRequest{
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
		return false, fmt.Errorf("%w: %s", appointment.ErrBusyPeriodsLoadFailed, err)
	}
	return len(res.Results) > 0, nil
}

func (r *RecordRepository) SaveRecord(ctx context.Context, rec *appointment.RecordEntity) error {
	start := notionapi.Date(entity.DateTimeToGoTime(rec.DateTimePeriod.Start))
	end := notionapi.Date(entity.DateTimeToGoTime(rec.DateTimePeriod.End))
	status, err := RecordStatusToNotion(*rec)
	if err != nil {
		return err
	}
	properties := notionapi.Properties{
		RecordTitle: notionapi.TitleProperty{
			Type:  notionapi.PropertyTypeTitle,
			Title: notion.ToRichText(rec.Customer().Name),
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
					ID: notionapi.PageID(rec.Customer().Id.String()),
				},
			},
		},
		RecordService: notionapi.RelationProperty{
			Type: notionapi.PropertyTypeRelation,
			Relation: []notionapi.Relation{
				{
					ID: notionapi.PageID(rec.Service().Id.String()),
				},
			},
		},
	}
	res, err := r.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: r.appointmentsDatabaseId,
		},
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", appointment.ErrAppointmentSaveFailed, err.Error())
	}
	rec.SetId(appointment.NewRecordId(string(res.ID)))
	return nil
}
