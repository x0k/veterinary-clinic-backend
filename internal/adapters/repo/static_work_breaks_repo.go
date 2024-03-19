package repo

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var allWorkBreaks = entity.WorkBreaks{
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

type StaticWorkBreaks struct{}

func NewStaticWorkBreaks() *StaticWorkBreaks {
	return &StaticWorkBreaks{}
}

func (s *StaticWorkBreaks) WorkBreaks(ctx context.Context) (entity.WorkBreaks, error) {
	return allWorkBreaks, nil
}
