package appointment_notion_repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

const appointmentRepositoryName = "appointment_notion_repository.AppointmentRepository"

type AppointmentRepository struct {
	log               *logger.Logger
	client            *notionapi.Client
	recordsDatabaseId notionapi.DatabaseID
}

func NewAppointment(
	log *logger.Logger,
	client *notionapi.Client,
	recordsDatabaseId notionapi.DatabaseID,
) *AppointmentRepository {
	return &AppointmentRepository{
		log:               log.With(slog.String("component", appointmentRepositoryName)),
		client:            client,
		recordsDatabaseId: recordsDatabaseId,
	}
}

func (r *AppointmentRepository) CreateAppointment(ctx context.Context, app *appointment.AppointmentAggregate) error {
	const op = appointmentRepositoryName + ".CreateAppointment"
	period := app.DateTimePeriod()
	start := notionapi.Date(entity.DateTimeToGoTime(period.Start))
	end := notionapi.Date(entity.DateTimeToGoTime(period.End))
	status, err := RecordStatusToNotion(app.Status(), app.IsArchived())
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
		// RecordCreatedAt: notionapi.CreatedTimeProperty{
		// 	Type:        notionapi.PropertyTypeCreatedTime,
		// 	CreatedTime: app.CreatedAt(),
		// },
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
	return app.SetId(appointment.NewRecordId(string(res.ID)))
}

func (s *AppointmentRepository) BusyPeriods(ctx context.Context, t time.Time) (appointment.BusyPeriods, error) {
	const op = appointmentRepositoryName + ".BusyPeriods"
	after := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	afterDate := notionapi.Date(after)
	beforeDate := notionapi.Date(after.AddDate(0, 0, 1))
	r, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: RecordDateTimePeriod,
				Date: &notionapi.DateFilterCondition{
					After: &afterDate,
				},
			},
			notionapi.PropertyFilter{
				Property: RecordDateTimePeriod,
				Date: &notionapi.DateFilterCondition{
					Before: &beforeDate,
				},
			},
			notionapi.PropertyFilter{
				Property: RecordState,
				Select: &notionapi.SelectFilterCondition{
					Equals: RecordAwaits,
				},
			},
		},
		Sorts: []notionapi.SortObject{
			{
				Property:  RecordDateTimePeriod,
				Direction: notionapi.SortOrderASC,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	periods := make([]entity.TimePeriod, 0, len(r.Results))
	for _, page := range r.Results {
		period, err := notion.DatePeriod(page.Properties, RecordDateTimePeriod)
		if err != nil {
			s.log.Error(ctx, "failed to parse record period", sl.Err(err))
			continue
		}
		periods = append(periods, entity.TimePeriod{
			Start: entity.GoTimeToTime(period.Start),
			End:   entity.GoTimeToTime(period.End),
		})
	}
	return periods, nil
}
