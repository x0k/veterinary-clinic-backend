package appointment

import "maps"

type AppointmentsState struct {
	appointments map[RecordId]AppointmentAggregate
}

func NewAppointmentsState(appointments map[RecordId]AppointmentAggregate) AppointmentsState {
	return AppointmentsState{
		appointments: appointments,
	}
}

func (s *AppointmentsState) Appointments() map[RecordId]AppointmentAggregate {
	return s.appointments
}

func (s *AppointmentsState) Reconcile(actualAppointments []AppointmentAggregate) []ChangedEvent {
	appsCopy := maps.Clone(s.appointments)
	changes := make([]ChangedEvent, 0, len(actualAppointments))
	for _, actualApp := range actualAppointments {
		s.appointments[actualApp.Id()] = actualApp
		oldApp, ok := appsCopy[actualApp.Id()]
		// created
		if !ok {
			changes = append(changes, ChangedEvent{
				ChangeType:  CreatedChangeType,
				Appointment: actualApp,
			})
			continue
		}
		if oldApp.Status() != actualApp.Status() {
			changes = append(changes, ChangedEvent{
				ChangeType:  StatusChangeType,
				Appointment: actualApp,
			})
		} else if oldApp.DateTimePeriod() != actualApp.DateTimePeriod() {
			changes = append(changes, ChangedEvent{
				ChangeType:  DateTimeChangeType,
				Appointment: actualApp,
			})
		}
		delete(appsCopy, actualApp.Id())
	}
	for _, app := range appsCopy {
		delete(s.appointments, app.Id())
		changes = append(changes, ChangedEvent{
			ChangeType:  RemovedChangeType,
			Appointment: app,
		})
	}
	return changes
}

func (s *AppointmentsState) AddAppointment(appointment AppointmentAggregate) {
	s.appointments[appointment.Id()] = appointment
}

func (s *AppointmentsState) RemoveAppointment(appointment AppointmentAggregate) {
	delete(s.appointments, appointment.Id())
}
