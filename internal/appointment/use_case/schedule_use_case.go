package appointment_use_case

type ScheduleUseCase[R any] struct {
	productionCalendarLoader ProductionCalendarLoader
	workingHoursLoader       WorkingHoursLoader
	busyPeriodsLoader        BusyPeriodsLoader
	workBreaksLoader         WorkBreaksLoader
}
