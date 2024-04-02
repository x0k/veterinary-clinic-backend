package appointment_notion

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

type AppointmentRepository struct {
	periodsMu              sync.Mutex
	periods                []entity.DateTimePeriod
	client                 *notionapi.Client
	appointmentsDatabaseId notionapi.DatabaseID
	servicesDatabaseId     notionapi.DatabaseID
}

func NewAppointment(
	client *notionapi.Client,
	appointmentsDatabaseId notionapi.DatabaseID,
	servicesDatabaseId notionapi.DatabaseID,
) *AppointmentRepository {
	return &AppointmentRepository{
		client:                 client,
		appointmentsDatabaseId: appointmentsDatabaseId,
		servicesDatabaseId:     servicesDatabaseId,
	}
}

func (r *AppointmentRepository) LockPeriod(ctx context.Context, period entity.DateTimePeriod) error {
	r.periodsMu.Lock()
	defer r.periodsMu.Unlock()
	for _, p := range r.periods {
		if entity.DateTimePeriodApi.IsValidPeriod(
			entity.DateTimePeriodApi.IntersectPeriods(p, period),
		) {
			return fmt.Errorf("%w: %s", appointment.ErrPeriodIsLocked, period)
		}
	}
	r.periods = append(r.periods, period)
	return nil
}

func (r *AppointmentRepository) UnLockPeriod(ctx context.Context, period entity.DateTimePeriod) error {
	r.periodsMu.Lock()
	defer r.periodsMu.Unlock()
	index := slices.Index(r.periods, period)
	if index == -1 {
		return nil
	}
	r.periods = slices.Delete(r.periods, index, index+1)
	return nil
}

func (r *AppointmentRepository) IsAppointmentPeriodBusy(ctx context.Context, period entity.DateTimePeriod) (bool, error) {
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
		return false, fmt.Errorf("%w: %s", appointment.ErrAppointmentsLoadFailed, err)
	}
	return len(res.Results) > 0, nil
}

func (r *AppointmentRepository) SaveAppointment(ctx context.Context, app *appointment.AppointmentAggregate) error {
	start := notionapi.Date(entity.DateTimeToGoTime(app.Record().DateTimePeriod.Start))
	end := notionapi.Date(entity.DateTimeToGoTime(app.Record().DateTimePeriod.End))
	status, err := RecordStatus(app.Record())
	if err != nil {
		return err
	}
	properties := notionapi.Properties{
		RecordTitle: notionapi.TitleProperty{
			Type:  notionapi.PropertyTypeTitle,
			Title: notion.ToRichText(app.Client().Name),
		},
		RecordService: notionapi.RelationProperty{
			Type: notionapi.PropertyTypeRelation,
			Relation: []notionapi.Relation{
				{
					ID: notionapi.PageID(app.Service().Id.String()),
				},
			},
		},
		RecordEmail: notionapi.EmailProperty{
			Type:  notionapi.PropertyTypeEmail,
			Email: app.Client().Email,
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
		RecordUserId: notionapi.RichTextProperty{
			Type:     notionapi.PropertyTypeRichText,
			RichText: notion.ToRichText(app.Client().Id.String()),
		},
	}
	if app.Client().PhoneNumber != "" {
		properties[RecordPhoneNumber] = notionapi.PhoneNumberProperty{
			Type:        notionapi.PropertyTypePhoneNumber,
			PhoneNumber: app.Client().PhoneNumber,
		}
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
	app.SetId(appointment.NewRecordId(string(res.ID)))
	return nil
}
