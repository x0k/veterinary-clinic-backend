package repo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

var ErrFailedToCreateRecord = errors.New("failed to create record")
var ErrLoadingBusyPeriodsFailed = errors.New("failed to load busy periods")

type NotionRecordsRepo struct {
	log                *logger.Logger
	recordsDatabaseId  notionapi.DatabaseID
	servicesDatabaseId notionapi.DatabaseID
	client             *notionapi.Client
	servicesCache      *containers.Expiable[[]entity.Service]
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
		servicesCache:      containers.NewExpiable[[]entity.Service](time.Hour),
	}
}

func (s *NotionRecordsRepo) Start(ctx context.Context) error {
	s.servicesCache.Start(ctx)
	return nil
}

func (s *NotionRecordsRepo) BusyPeriods(ctx context.Context, t time.Time) (entity.BusyPeriods, error) {
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
	periods := make([]entity.TimePeriod, 0, len(r.Results))
	for _, page := range r.Results {
		if period := DateTimePeriodFromRecord(page.Properties); period != nil {
			periods = append(periods, entity.TimePeriod{
				Start: period.Start.Time,
				End:   period.End.Time,
			})
		}
	}
	return periods, nil
}

func (s *NotionRecordsRepo) Services(ctx context.Context) ([]entity.Service, error) {
	return s.servicesCache.Load(func() ([]entity.Service, error) {
		r, err := s.client.Database.Query(ctx, s.servicesDatabaseId, nil)
		if err != nil {
			return nil, err
		}
		services := make([]entity.Service, 0, len(r.Results))
		for _, result := range r.Results {
			services = append(services, Service(result))
		}
		return services, nil
	})
}

func (s *NotionRecordsRepo) Service(ctx context.Context, serviceId entity.ServiceId) (entity.Service, error) {
	r, err := s.client.Page.Get(ctx, notionapi.PageID(serviceId))
	if err != nil {
		return entity.Service{}, err
	}
	if r == nil {
		return entity.Service{}, usecase.ErrNotFound
	}
	return Service(*r), nil
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
	service, err := s.Service(ctx, entity.ServiceId((relations[0].ID)))
	if err != nil {
		return entity.Record{}, err
	}
	if rec := ActualRecord(page, &userId, service); rec != nil {
		return *rec, nil
	}
	return entity.Record{}, ErrFailedToCreateRecord
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

func (s *NotionRecordsRepo) LoadActualRecords(ctx context.Context, now time.Time) ([]entity.Record, error) {
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
	servicesMap := make(map[entity.ServiceId]entity.Service, len(services))
	for _, service := range services {
		servicesMap[service.Id] = service
	}
	records := make([]entity.Record, 0, len(res.Results))
	errs := make([]error, 0, len(res.Results))
	for _, page := range res.Results {
		relations := Relations(page.Properties, RecordService)
		if len(relations) == 0 {
			errs = append(errs, adapters.ErrInvalidRecord)
			continue
		}
		service, ok := servicesMap[entity.ServiceId(relations[0].ID)]
		if !ok {
			errs = append(errs, adapters.ErrInvalidRecord)
			continue
		}
		rec, err := PrivateActualRecord(page, service)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		records = append(records, rec)
	}
	return records, errors.Join(errs...)
}
