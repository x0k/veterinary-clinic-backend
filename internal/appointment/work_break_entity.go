package appointment

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrFailedToCompileMatchExpression = errors.New("failed to compile match expression")

type WorkBreakId string

func NewWorkBreakId(id string) WorkBreakId {
	return WorkBreakId(id)
}

type WorkBreak struct {
	Id              WorkBreakId
	Title           string
	MatchExpression string
	Period          entity.TimePeriod
}

func NewWorkBreak(
	id WorkBreakId,
	title string,
	matchExpression string,
	period entity.TimePeriod,
) WorkBreak {
	return WorkBreak{
		Id:              id,
		Title:           title,
		MatchExpression: matchExpression,
		Period:          period,
	}
}

type WorkBreaks []WorkBreak

const date_format = "2006-01-02T15:04:05"

func (workBreaks WorkBreaks) ForDay(
	t time.Time,
) (WorkBreaks, error) {
	dateString := fmt.Sprintf("%d %s", t.Weekday(), t.Format(date_format))
	breaks := make(WorkBreaks, 0, len(workBreaks))
	for _, wb := range workBreaks {
		expr, err := regexp.Compile(wb.MatchExpression)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFailedToCompileMatchExpression, wb.MatchExpression)
		}
		if expr.MatchString(dateString) {
			breaks = append(breaks, wb)
		}
	}
	return breaks, nil
}
