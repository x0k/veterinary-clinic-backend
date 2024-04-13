package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var staticWorkBreaks = shared.WorkBreaks{
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
	breaksCache      *containers.Expiable[shared.WorkBreaks]
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
		breaksCache:      containers.NewExpiable[shared.WorkBreaks](time.Hour),
	}
}

func (s *NotionWorkBreaks) Start(ctx context.Context) error {
	s.breaksCache.Start(ctx)
	return nil
}

func (s *NotionWorkBreaks) WorkBreaks(ctx context.Context) (shared.WorkBreaks, error) {
	return s.breaksCache.Load(func() (shared.WorkBreaks, error) {
		r, err := s.client.Database.Query(ctx, s.breaksDatabaseId, nil)
		if err != nil {
			return nil, err
		}
		workBreaks := make(shared.WorkBreaks, len(staticWorkBreaks), len(staticWorkBreaks)+len(r.Results))
		copy(workBreaks, staticWorkBreaks)
		for _, result := range r.Results {
			workBreak, err := WorkBreak(result)
			if err != nil {
				s.log.Error(ctx, "failed to parse work break", sl.Err(err))
				continue
			}
			workBreaks = append(workBreaks, workBreak)
		}
		return workBreaks, nil
	})
}
