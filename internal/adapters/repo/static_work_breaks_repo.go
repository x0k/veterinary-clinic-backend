package repo

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var wbCalculator = entity.NewWorkBreaksCalculator(
	[]entity.WorkBreak{
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
					Minutes: 0,
				},
			},
		},
	},
)

type Static struct{}

func NewStatic() *Static {
	return &Static{}
}

func (s *Static) WorkBreaks(ctx context.Context, t time.Time) (entity.WorkBreaks, error) {
	return wbCalculator.Calculate(t)
}
