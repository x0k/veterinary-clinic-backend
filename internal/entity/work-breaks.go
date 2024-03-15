package entity

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var ErrFailedToCompileMatchExpression = errors.New("failed to compile match expression")

type WorkBreakId string

type WorkBreak struct {
	Id              WorkBreakId
	Title           string
	MatchExpression string
	Period          TimePeriod
}

type WorkBreaks []WorkBreak

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
) (WorkBreaks, error) {
	dateString := fmt.Sprintf("%d %s", t.Weekday(), t.Format(date_format))
	breaks := make(WorkBreaks, 0, len(c.workBreaks))
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
