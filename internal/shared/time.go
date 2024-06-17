package shared

import (
	"errors"
	"time"
)

var ErrInvalidWeekday = errors.New("invalid weekday")

func NewWeekday(d int) (time.Weekday, error) {
	if d < 0 || d > 6 {
		return 0, ErrInvalidWeekday
	}
	return time.Weekday(d), nil
}
