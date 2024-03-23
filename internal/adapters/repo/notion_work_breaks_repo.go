package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

var staticWorkBreaks = entity.WorkBreaks{
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

type NotionWorkBreaks struct {
	log              *logger.Logger
	breaksDatabaseId notionapi.DatabaseID
	client           *notionapi.Client
	breaksCache      *containers.Expiable[entity.WorkBreaks]
}

func NewNotionWorkBreaks(
	client *notionapi.Client,
	log *logger.Logger,
	breaksDatabaseId notionapi.DatabaseID,
) *NotionWorkBreaks {
	return &NotionWorkBreaks{
		log:              log.With(slog.String("component", "adapters.repo.NotionBreaksRepo")),
		client:           client,
		breaksDatabaseId: breaksDatabaseId,
		breaksCache:      containers.NewExpiable[entity.WorkBreaks](time.Hour),
	}
}

func (s *NotionWorkBreaks) Start(ctx context.Context) error {
	s.breaksCache.Start(ctx)
	return nil
}

func (s *NotionWorkBreaks) WorkBreaks(ctx context.Context) (entity.WorkBreaks, error) {
	return s.breaksCache.Load(func() (entity.WorkBreaks, error) {
		r, err := s.client.Database.Query(ctx, s.breaksDatabaseId, nil)
		if err != nil {
			return nil, err
		}
		workBreaks := make(entity.WorkBreaks, len(staticWorkBreaks), len(staticWorkBreaks)+len(r.Results))
		copy(workBreaks, staticWorkBreaks)
		for _, result := range r.Results {
			if workBreak := WorkBreak(result); workBreak != nil {
				workBreaks = append(workBreaks, *workBreak)
			}
		}
		return workBreaks, nil
	})
}
