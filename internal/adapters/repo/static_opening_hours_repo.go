package repo

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var weekdayTimePeriod = entity.TimePeriod{
	Start: entity.Time{
		Hours:   9,
		Minutes: 30,
	},
	End: entity.Time{
		Hours:   17,
		Minutes: 0,
	},
}
var saturdayTimePeriod = entity.TimePeriod{
	Start: weekdayTimePeriod.Start,
	End: entity.Time{
		Hours:   13,
		Minutes: 0,
	},
}
var openingHours = entity.OpeningHours{
	1: weekdayTimePeriod,
	2: weekdayTimePeriod,
	3: weekdayTimePeriod,
	4: weekdayTimePeriod,
	5: weekdayTimePeriod,
	6: saturdayTimePeriod,
}

type StaticOpeningHours struct{}

func NewStaticOpeningHours() *StaticOpeningHours {
	return &StaticOpeningHours{}
}

func (s *StaticOpeningHours) OpeningHours(t time.Time) entity.OpeningHours {
	return openingHours
}
