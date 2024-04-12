package appointment_notion_repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/notion"
)

const appointmentRepositoryName = "appointment_notion_repository.AppointmentRepository"

type AppointmentRepository struct {
	log                *logger.Logger
	client             *notionapi.Client
	recordsDatabaseId  notionapi.DatabaseID
	servicesDatabaseId notionapi.DatabaseID
	servicesCache      *containers.Expiable[[]appointment.ServiceEntity]
}

func NewAppointment(
	log *logger.Logger,
	client *notionapi.Client,
	recordsDatabaseId notionapi.DatabaseID,
	servicesDatabaseId notionapi.DatabaseID,
) *AppointmentRepository {
	return &AppointmentRepository{
		log:                log.With(slog.String("component", appointmentRepositoryName)),
		client:             client,
		recordsDatabaseId:  recordsDatabaseId,
		servicesDatabaseId: servicesDatabaseId,
		servicesCache:      containers.NewExpiable[[]appointment.ServiceEntity](time.Hour),
	}
}

func (r *AppointmentRepository) Name() string {
	return appointmentRepositoryName
}

func (r *AppointmentRepository) Start(ctx context.Context) error {
	r.servicesCache.Start(ctx)
	return nil
}

func (s *AppointmentRepository) Services(ctx context.Context) ([]appointment.ServiceEntity, error) {
	const op = appointmentRepositoryName + ".Services"
	return s.servicesCache.Load(func() ([]appointment.ServiceEntity, error) {
		r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		services := make([]appointment.ServiceEntity, 0, len(r.Results))
		for _, result := range r.Results {
			services = append(services, NotionToService(result))
		}
		return services, nil
	})
}

func (s *AppointmentRepository) Service(ctx context.Context, serviceId appointment.ServiceId) (appointment.ServiceEntity, error) {
	const op = appointmentRepositoryName + ".Service"
	res, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil {
		return appointment.ServiceEntity{}, fmt.Errorf("%s: %w", op, entity.ErrNotFound)
	}
	return NotionToService(*res), nil
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

func (s *AppointmentRepository) CustomerActiveAppointment(
	ctx context.Context,
	customer appointment.CustomerEntity,
) (appointment.AppointmentAggregate, error) {
	const op = appointmentRepositoryName + ".CustomerActiveAppointment"
	res, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: RecordCustomer,
				Relation: &notionapi.RelationFilterCondition{
					Contains: customer.Id.String(),
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
		return appointment.AppointmentAggregate{}, fmt.Errorf("%s: %w", op, err)
	}
	if res == nil || len(res.Results) == 0 {
		return appointment.AppointmentAggregate{}, fmt.Errorf("%s: %w", op, entity.ErrNotFound)
	}
	record, err := NotionToRecord(res.Results[0])
	if err != nil {
		return appointment.AppointmentAggregate{}, fmt.Errorf("%s: %w", op, err)
	}
	service, err := s.Service(ctx, record.ServiceId)
	if err != nil {
		return appointment.AppointmentAggregate{}, fmt.Errorf("%s: %w", op, err)
	}
	return appointment.NewAppointmentAggregate(record, service, customer)
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
