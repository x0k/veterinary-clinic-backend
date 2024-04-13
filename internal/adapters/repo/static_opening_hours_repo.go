package repo

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var weekdayTimePeriod = shared.TimePeriod{
	Start: shared.Time{
		Hours:   9,
		Minutes: 30,
	},
	End: shared.Time{
		Hours:   17,
		Minutes: 0,
	},
}
var saturdayTimePeriod = shared.TimePeriod{
	Start: weekdayTimePeriod.Start,
	End: shared.Time{
		Hours:   13,
		Minutes: 0,
	},
}
var openingHours = shared.OpeningHours{
	1: weekdayTimePeriod,
	2: weekdayTimePeriod,
	3: weekdayTimePeriod,
	4: weekdayTimePeriod,
	5: weekdayTimePeriod,
	6: saturdayTimePeriod,
}

type StaticOpeningHoursRepo struct{}

func NewStaticOpeningHoursRepo() *StaticOpeningHoursRepo {
	return &StaticOpeningHoursRepo{}
}

func (r *StaticOpeningHoursRepo) OpeningHours(ctx context.Context) (shared.OpeningHours, error) {
	return openingHours, nil
}
