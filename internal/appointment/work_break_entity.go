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

type WorkBreak struct {
	Id              WorkBreakId
	Title           string
	MatchExpression string
	Period          entity.TimePeriod
}

type WorkBreaks []WorkBreak

type CalculatedWorkBreaks []WorkBreak

type WorkBreaksCalculator struct {
	workBreaks WorkBreaks
}

func NewWorkBreaksCalculator(
	workBreaks WorkBreaks,
) *WorkBreaksCalculator {
	return &WorkBreaksCalculator{
		workBreaks: workBreaks,
	}
}

const date_format = "2006-01-02T15:04:05"

func (c *WorkBreaksCalculator) Calculate(
	t time.Time,
) (CalculatedWorkBreaks, error) {
	dateString := fmt.Sprintf("%d %s", t.Weekday(), t.Format(date_format))
	breaks := make(CalculatedWorkBreaks, 0, len(c.workBreaks))
	for _, wb := range c.workBreaks {
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

func CalculateWorkBreaks(workBreaks WorkBreaks, date time.Time) (CalculatedWorkBreaks, error) {
	return NewWorkBreaksCalculator(workBreaks).Calculate(date)
}
