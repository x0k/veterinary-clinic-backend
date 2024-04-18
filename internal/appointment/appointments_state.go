package appointment

import "maps"

type AppointmentsState struct {
	appointments map[RecordId]RecordEntity
}

func NewAppointmentsState(appointments map[RecordId]RecordEntity) AppointmentsState {
	return AppointmentsState{
		appointments: appointments,
	}
}

func (s *AppointmentsState) Appointments() map[RecordId]RecordEntity {
	return s.appointments
}

func (s *AppointmentsState) Reconcile(actualAppointments []RecordEntity) []ChangedEvent {
	appsCopy := maps.Clone(s.appointments)
	changes := make([]ChangedEvent, 0, len(actualAppointments))
	for _, actualApp := range actualAppointments {
		s.appointments[actualApp.Id] = actualApp
		oldApp, ok := appsCopy[actualApp.Id]
		// created
		if !ok {
			changes = append(changes, NewChanged(
				CreatedChangeType,
				actualApp,
			))
			continue
		}
		if oldApp.Status != actualApp.Status {
			changes = append(changes, NewChanged(
				StatusChangeType,
				actualApp,
			))
		} else if oldApp.DateTimePeriod != actualApp.DateTimePeriod {
			changes = append(changes, NewChanged(
				DateTimeChangeType,
				actualApp,
			))
		}
		delete(appsCopy, actualApp.Id)
	}
	for _, app := range appsCopy {
		delete(s.appointments, app.Id)
		changes = append(changes, NewChanged(
			RemovedChangeType,
			app,
		))
	}
	return changes
}

func (s *AppointmentsState) AddAppointment(appointment RecordEntity) {
	s.appointments[appointment.Id] = appointment
}

func (s *AppointmentsState) RemoveAppointment(appointment RecordEntity) {
	delete(s.appointments, appointment.Id)
}
