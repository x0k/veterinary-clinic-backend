package appointment_notion_repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const workBreaksRepositoryName = "appointment_notion_repository.WorkBreaksRepository"

var staticWorkBreaks = appointment.WorkBreaks{
	{
		Id:              "lunch",
		MatchExpression: `^[1-5]`,
		Title:           "Перерыв на обед",
		Period: entity.TimePeriod{
			Start: entity.Time{
				Hours:   12,
				Minutes: 30,
			},
			End: entity.Time{
				Hours:   13,
				Minutes: 30,
			},
		},
	},
}

type WorkBreaksRepository struct {
	log              *logger.Logger
	breaksDatabaseId notionapi.DatabaseID
	client           *notionapi.Client
	breaksCache      *containers.Expiable[appointment.WorkBreaks]
}

func NewWorkBreaks(
	log *logger.Logger,
	client *notionapi.Client,
	breaksDatabaseId notionapi.DatabaseID,
) *WorkBreaksRepository {
	return &WorkBreaksRepository{
		log:              log.With(slog.String("component", workBreaksRepositoryName)),
		client:           client,
		breaksDatabaseId: breaksDatabaseId,
		breaksCache:      containers.NewExpiable[appointment.WorkBreaks](time.Hour),
	}
}

func (s *WorkBreaksRepository) Name() string {
	return workBreaksRepositoryName
}

func (s *WorkBreaksRepository) Start(ctx context.Context) error {
	s.breaksCache.Start(ctx)
	return nil
}

func (s *WorkBreaksRepository) WorkBreaks(ctx context.Context) (appointment.WorkBreaks, error) {
	return s.breaksCache.Load(func() (appointment.WorkBreaks, error) {
		r, err := s.client.Database.Query(ctx, s.breaksDatabaseId, nil)
		if err != nil {
			return nil, err
		}
		workBreaks := make(appointment.WorkBreaks, len(staticWorkBreaks), len(staticWorkBreaks)+len(r.Results))
		copy(workBreaks, staticWorkBreaks)
		for _, result := range r.Results {
			workBreak, err := NotionToWorkBreak(result)
			if err != nil {
				s.log.Error(ctx, "failed to parse work break", sl.Err(err))
				continue
			}
			workBreaks = append(workBreaks, workBreak)
		}
		return workBreaks, nil
	})
}
