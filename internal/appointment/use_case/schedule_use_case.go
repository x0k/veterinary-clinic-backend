package appointment_use_case

import "github.com/x0k/veterinary-clinic-backend/internal/appointment"

type ScheduleUseCase[R any] struct {
	productionCalendarLoader appointment.ProductionCalendarLoader
	workingHoursLoader       appointment.WorkingHoursLoader
	busyPeriodsLoader        appointment.BusyPeriodsLoader
	workBreaksLoader         appointment.WorkBreaksLoader
}
