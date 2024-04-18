package appointment

import (
	"context"
	"sync"
	"time"
)

type TrackingService struct {
	appointmentsLoader ActualAppointmentsLoader
	stateMu            sync.Mutex
	stateLoader        AppointmentsStateLoader
	stateSaver         AppointmentsStateSaver
}

func NewTracking(
	appointmentsLoader ActualAppointmentsLoader,
	stateLoader AppointmentsStateLoader,
	stateSaver AppointmentsStateSaver,
) *TrackingService {
	return &TrackingService{
		appointmentsLoader: appointmentsLoader,
		stateLoader:        stateLoader,
		stateSaver:         stateSaver,
	}
}

func (s *TrackingService) state(
	ctx context.Context,
	mutate func(*AppointmentsState),
) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	state, err := s.stateLoader(ctx)
	if err != nil {
		return err
	}
	mutate(&state)
	return s.stateSaver(ctx, state)
}

func (s *TrackingService) DetectChanges(
	ctx context.Context,
	now time.Time,
) ([]ChangedEvent, error) {
	actualAppointments, err := s.appointmentsLoader(ctx, now)
	if err != nil {
		return nil, err
	}
	var changes []ChangedEvent
	err = s.state(ctx, func(state *AppointmentsState) {
		changes = state.Reconcile(actualAppointments)
	})
	return changes, err
}

func (s *TrackingService) AddAppointment(
	ctx context.Context,
	appointment RecordEntity,
) error {
	return s.state(ctx, func(state *AppointmentsState) {
		state.AddAppointment(appointment)
	})
}

func (s *TrackingService) RemoveAppointment(
	ctx context.Context,
	appointment RecordEntity,
) error {
	return s.state(ctx, func(state *AppointmentsState) {
		state.RemoveAppointment(appointment)
	})
}
