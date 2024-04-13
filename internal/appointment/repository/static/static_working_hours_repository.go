package appointment_static_repository

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
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
var workingHours = appointment.WorkingHours{
	1: weekdayTimePeriod,
	2: weekdayTimePeriod,
	3: weekdayTimePeriod,
	4: weekdayTimePeriod,
	5: weekdayTimePeriod,
	6: saturdayTimePeriod,
}

type WorkingHoursRepository struct{}

func NewWorkingHoursRepository() *WorkingHoursRepository {
	return &WorkingHoursRepository{}
}

func (r *WorkingHoursRepository) WorkingHours(ctx context.Context) (appointment.WorkingHours, error) {
	return workingHours, nil
}
