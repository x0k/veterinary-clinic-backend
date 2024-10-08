package appointment_notion_repository

import (
	"context"

	"github.com/jomei/notionapi"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const workBreaksRepositoryName = "appointment_notion_repository.WorkBreaksRepository"

var staticWorkBreaks = appointment.WorkBreaks{
	{
		Id:              "lunch",
		MatchExpression: `^[1-5]`,
		Title:           "Перерыв на обед",
		Period: shared.TimePeriod{
			Start: shared.Time{
				Hours:   12,
				Minutes: 30,
			},
			End: shared.Time{
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
}

func NewWorkBreaks(
	log *logger.Logger,
	client *notionapi.Client,
	breaksDatabaseId notionapi.DatabaseID,
) *WorkBreaksRepository {
	return &WorkBreaksRepository{
		log:              log,
		client:           client,
		breaksDatabaseId: breaksDatabaseId,
	}
}

func (s *WorkBreaksRepository) WorkBreaks(ctx context.Context) (appointment.WorkBreaks, error) {
	const op = workBreaksRepositoryName + ".WorkBreaks"
	r, err := s.client.Database.Query(ctx, s.breaksDatabaseId, nil)
	if err != nil {
		return nil, err
	}
	workBreaks := make(appointment.WorkBreaks, len(staticWorkBreaks), len(staticWorkBreaks)+len(r.Results))
	copy(workBreaks, staticWorkBreaks)
	for _, result := range r.Results {
		workBreak, err := NotionToWorkBreak(result)
		if err != nil {
			s.log.Error(ctx, "failed to parse work break", sl.Op(op), sl.Err(err))
			continue
		}
		workBreaks = append(workBreaks, workBreak)
	}
	return workBreaks, nil
}
