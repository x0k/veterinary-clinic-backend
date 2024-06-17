package appointment_notion_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
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
		log:               log,
		client:            client,
		recordsDatabaseId: recordsDatabaseId,
	}
}

func (r *AppointmentRepository) CreateAppointment(ctx context.Context, app *appointment.RecordEntity) error {
	const op = appointmentRepositoryName + ".CreateAppointment"
	period := app.DateTimePeriod
	start := notionapi.Date(shared.DateTimeToUTCTime(period.Start).Time)
	end := notionapi.Date(shared.DateTimeToUTCTime(period.End).Time)
	status, err := RecordStatusToNotion(app.Status, app.IsArchived)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	properties := notionapi.Properties{
		RecordTitle: notionapi.TitleProperty{
			Type:  notionapi.PropertyTypeTitle,
			Title: notion.ToRichText(app.Title),
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
					ID: notionapi.PageID(app.CustomerId.String()),
				},
			},
		},
		RecordService: notionapi.RelationProperty{
			Type: notionapi.PropertyTypeRelation,
			Relation: []notionapi.Relation{
				{
					ID: notionapi.PageID(app.ServiceId.String()),
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
	app.SetCreatedAt(
		notion.CreatedTime(res.Properties, RecordCreatedAt),
	)
	return app.SetId(appointment.NewRecordId(res.ID.String()))
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
	periods := make([]shared.TimePeriod, 0, len(r.Results))
	for _, page := range r.Results {
		period, err := notion.DatePeriod(page.Properties, RecordDateTimePeriod)
		if err != nil {
			s.log.Error(ctx, "failed to parse record period", sl.Op(op), sl.Err(err))
			continue
		}
		periods = append(periods, shared.TimePeriod{
			Start: shared.UTCTimeToTime(shared.NewUTCTime(period.Start)),
			End:   shared.UTCTimeToTime(shared.NewUTCTime(period.End)),
		})
	}
	return periods, nil
}

func (s *AppointmentRepository) CustomerActiveAppointment(
	ctx context.Context,
	customerId appointment.CustomerId,
) (appointment.RecordEntity, error) {
	const op = appointmentRepositoryName + ".CustomerActiveAppointment"
	res, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: RecordCustomer,
				Relation: &notionapi.RelationFilterCondition{
					Contains: customerId.String(),
				},
			},
			notionapi.OrCompoundFilter{
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordAwaits,
					},
				},
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordDone,
					},
				},
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordNotAppear,
					},
				},
			},
		},
	})
	if err != nil {
		return appointment.RecordEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil || len(res.Results) == 0 {
		return appointment.RecordEntity{}, fmt.Errorf("%s: %w", op, shared.ErrNotFound)
	}
	return NotionToRecord(res.Results[0])
}

func (s *AppointmentRepository) RemoveAppointment(ctx context.Context, recordId appointment.RecordId) error {
	const op = appointmentRepositoryName + ".RemoveAppointment"
	_, err := s.client.Page.Update(ctx, notionapi.PageID(recordId.String()), &notionapi.PageUpdateRequest{
		Archived:   true,
		Properties: notionapi.Properties{},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *AppointmentRepository) ActualAppointments(
	ctx context.Context,
	now time.Time,
) ([]appointment.RecordEntity, error) {
	const op = appointmentRepositoryName + ".ActualAppointments"
	after := notionapi.Date(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
	recordsRes, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: RecordDateTimePeriod,
				Date: &notionapi.DateFilterCondition{
					After: &after,
				},
			},
			notionapi.OrCompoundFilter{
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordAwaits,
					},
				},
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordDone,
					},
				},
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordNotAppear,
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
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if recordsRes == nil || len(recordsRes.Results) == 0 {
		return nil, nil
	}
	records := make([]appointment.RecordEntity, 0, len(recordsRes.Results))
	for _, result := range recordsRes.Results {
		record, err := NotionToRecord(result)
		if err != nil {
			s.log.Error(ctx, "failed to convert record", sl.Op(op), sl.Err(err))
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func (r *AppointmentRepository) ArchiveRecords(ctx context.Context) error {
	res, err := r.client.Database.Query(ctx, r.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.OrCompoundFilter{
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordDone,
					},
				},
				notionapi.PropertyFilter{
					Property: RecordState,
					Select: &notionapi.SelectFilterCondition{
						Equals: RecordNotAppear,
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
	})
	if err != nil {
		return err
	}
	errs := make([]error, 0, len(res.Results))
	for _, page := range res.Results {
		status, _, err := NotionToRecordStatus(notion.Select(page.Properties, RecordState))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		newState := RecordDoneArchived
		if status == RecordNotAppear {
			newState = RecordNotAppearArchived
		}
		if _, err = r.client.Page.Update(ctx, notionapi.PageID(page.ID), &notionapi.PageUpdateRequest{
			Properties: notionapi.Properties{
				RecordState: notionapi.SelectProperty{
					Select: notionapi.Option{Name: newState},
				},
			},
		}); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
