package repo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

var ErrFailedToCreateRecord = errors.New("failed to create record")
var ErrLoadingBusyPeriodsFailed = errors.New("failed to load busy periods")

type NotionRecordsRepo struct {
	log                *logger.Logger
	recordsDatabaseId  notionapi.DatabaseID
	servicesDatabaseId notionapi.DatabaseID
	client             *notionapi.Client
	servicesCache      *containers.Expiable[[]shared.Service]
}

func NewNotionRecords(
	client *notionapi.Client,
	log *logger.Logger,
	recordsDatabaseId notionapi.DatabaseID,
	servicesDatabaseId notionapi.DatabaseID,
) *NotionRecordsRepo {
	return &NotionRecordsRepo{
		log:                log.With(slog.String("component", "adapters.repo.NotionRecordsRepo")),
		client:             client,
		recordsDatabaseId:  recordsDatabaseId,
		servicesDatabaseId: servicesDatabaseId,
		servicesCache:      containers.NewExpiable[[]shared.Service](time.Hour),
	}
}

func (s *NotionRecordsRepo) Start(ctx context.Context) error {
	s.servicesCache.Start(ctx)
	return nil
}

func (s *NotionRecordsRepo) BusyPeriods(ctx context.Context, t time.Time) (shared.BusyPeriods, error) {
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
		s.log.Error(ctx, "failed to load busy periods", sl.Err(err))
		return nil, ErrLoadingBusyPeriodsFailed
	}
	periods := make([]shared.TimePeriod, 0, len(r.Results))
	for _, page := range r.Results {
		period, err := DateTimePeriod(page.Properties, RecordDateTimePeriod)
		if err != nil {
			s.log.Error(ctx, "failed to parse record period", sl.Err(err))
			continue
		}
		periods = append(periods, shared.TimePeriod{
			Start: period.Start.Time,
			End:   period.End.Time,
		})
	}
	return periods, nil
}

func (s *NotionRecordsRepo) Services(ctx context.Context) ([]shared.Service, error) {
	return s.servicesCache.Load(func() ([]shared.Service, error) {
		r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
		if err != nil {
			return nil, err
		}
		services := make([]shared.Service, 0, len(r.Results))
		for _, result := range r.Results {
			services = append(services, Service(result))
		}
		return services, nil
	})
}

func (s *NotionRecordsRepo) Service(ctx context.Context, serviceId shared.ServiceId) (shared.Service, error) {
	r, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return shared.Service{}, err
	}
	if r == nil {
		return shared.Service{}, usecase.ErrNotFound
	}
	return Service(*r), nil
}

func (s *NotionRecordsRepo) Create(
	ctx context.Context,
	user shared.User,
	service shared.Service,
	appointmentDateTime time.Time,
) (shared.Record, error) {
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
		return shared.Record{}, err
	}
	if res == nil {
		return shared.Record{}, ErrFailedToCreateRecord
	}
	return Record(*res, service)
}

func (s *NotionRecordsRepo) Remove(ctx context.Context, recordId shared.RecordId) error {
	_, err := s.client.Page.Update(ctx, notionapi.PageID(recordId), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{},
		Archived:   true,
	})
	return err
}

func (s *NotionRecordsRepo) RecordByUserId(ctx context.Context, userId shared.UserId) (shared.Record, error) {
	res, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
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
			notionapi.PropertyFilter{
				Property: RecordUserId,
				RichText: &notionapi.TextFilterCondition{
					Equals: string(userId),
				},
			},
		},
	})
	if err != nil {
		return shared.Record{}, err
	}
	if len(res.Results) == 0 {
		return shared.Record{}, usecase.ErrNotFound
	}
	page := res.Results[0]
	relations := Relations(page.Properties, RecordService)
	if len(relations) == 0 {
		return shared.Record{}, adapters.ErrInvalidRecord
	}
	service, err := s.Service(ctx, shared.ServiceId((relations[0].ID)))
	if err != nil {
		return shared.Record{}, err
	}
	return Record(page, service)
}

func (s *NotionRecordsRepo) LoadActualRecords(ctx context.Context, now time.Time) ([]shared.Record, error) {
	after := notionapi.Date(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
	res, err := s.client.Database.Query(ctx, s.recordsDatabaseId, &notionapi.DatabaseQueryRequest{
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
		return nil, err
	}
	if len(res.Results) == 0 {
		return nil, nil
	}
	services, err := s.Services(ctx)
	if err != nil {
		return nil, err
	}
	servicesMap := make(map[shared.ServiceId]shared.Service, len(services))
	for _, service := range services {
		servicesMap[service.Id] = service
	}
	records := make([]shared.Record, 0, len(res.Results))
	errs := make([]error, 0, len(res.Results))
	for _, page := range res.Results {
		relations := Relations(page.Properties, RecordService)
		if len(relations) == 0 {
			errs = append(errs, adapters.ErrInvalidRecord)
			continue
		}
		service, ok := servicesMap[shared.ServiceId(relations[0].ID)]
		if !ok {
			errs = append(errs, adapters.ErrInvalidRecord)
			continue
		}
		rec, err := Record(page, service)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		records = append(records, rec)
	}
	return records, errors.Join(errs...)
}

func (r *NotionRecordsRepo) ArchiveRecords(ctx context.Context) error {
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
		status, err := RecordStatus(page.Properties)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		newState := RecordDoneArchived
		if status == shared.RecordNotAppear {
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
