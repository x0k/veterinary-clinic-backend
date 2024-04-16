package appointment

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ActualAppointmentsState interface {
	Reconcile(actualAppointments []AppointmentAggregate) []ChangedEvent
	AddAppointment(appointment AppointmentAggregate)
	RemoveAppointment(appointment AppointmentAggregate)
}

type ChangesDetectorService struct {
	actualAppointmentsLoader shared.Loader[[]AppointmentAggregate]
	stateLoader              shared.Loader[ActualAppointmentsState]
	stateSaver               shared.Saver[ActualAppointmentsState]
}

func NewChangesDetectorService(
	stateLoader shared.Loader[ActualAppointmentsState],
	stateSaver shared.Saver[ActualAppointmentsState],
) *ChangesDetectorService {
	return &ChangesDetectorService{
		stateLoader: stateLoader,
		stateSaver:  stateSaver,
	}
}

func (s *ChangesDetectorService) state(
	ctx context.Context,
	mutate func(ActualAppointmentsState) error,
) error {
	state, err := s.stateLoader(ctx)
	if err != nil {
		return err
	}
	if err := mutate(state); err != nil {
		return err
	}
	return s.stateSaver(ctx, state)
}

func (s *ChangesDetectorService) DetectChanges(
	ctx context.Context,
) ([]ChangedEvent, error) {
	actualAppointments, err := s.actualAppointmentsLoader(ctx)
	if err != nil {
		return nil, err
	}
	var changes []ChangedEvent
	err = s.state(ctx, func(state ActualAppointmentsState) error {
		changes = state.Reconcile(actualAppointments)
		return nil
	})
	return changes, err
}
