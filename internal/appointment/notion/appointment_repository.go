package appointment_notion

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

type Appointment struct {
	log                    *logger.Logger
	periodsMu              sync.Mutex
	periods                []entity.DateTimePeriod
	client                 *notionapi.Client
	servicesCache          *containers.Expiable[[]appointment.Service]
	appointmentsDatabaseId notionapi.DatabaseID
	servicesDatabaseId     notionapi.DatabaseID
}

func NewAppointment(
	log *logger.Logger,
	client *notionapi.Client,
	appointmentsDatabaseId notionapi.DatabaseID,
	servicesDatabaseId notionapi.DatabaseID,
) *Appointment {
	return &Appointment{
		log:                    log.With(slog.String("component", "appointment_notion.Appointment")),
		client:                 client,
		appointmentsDatabaseId: appointmentsDatabaseId,
		servicesDatabaseId:     servicesDatabaseId,
		servicesCache:          containers.NewExpiable[[]appointment.Service](time.Hour),
	}
}

func (r *Appointment) LockPeriod(ctx context.Context, period entity.DateTimePeriod) error {
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

func (r *Appointment) UnLockPeriod(ctx context.Context, period entity.DateTimePeriod) error {
	r.periodsMu.Lock()
	defer r.periodsMu.Unlock()
	index := slices.Index(r.periods, period)
	if index == -1 {
		return nil
	}
	r.periods = slices.Delete(r.periods, index, index+1)
	return nil
}

func (r *Appointment) IsAppointmentPeriodBusy(ctx context.Context, period entity.DateTimePeriod) (bool, error) {
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
		r.log.Error(ctx, "failed to query database", sl.Err(err))
		return false, appointment.ErrAppointmentsLoadFailed
	}
	return len(res.Results) > 0, nil
}

func (r *Appointment) SaveAppointment(ctx context.Context, app *appointment.Appointment) error {
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
		r.log.Error(ctx, "failed to create page", sl.Err(err))
		return appointment.ErrAppointmentSaveFailed
	}
	app.SetId(appointment.NewRecordId(string(res.ID)))
	return nil
}
