package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type EventType int

const (
	CreatedEventType EventType = iota
	CanceledEventType
	ChangedEventType
)

type Event pubsub.Event[EventType]
type Publisher pubsub.Publisher[EventType]

type CreatedEvent struct {
	Record   RecordEntity
	Customer CustomerEntity
	Service  ServiceEntity
}

func NewCreated(
	appointment RecordEntity,
	customer CustomerEntity,
	service ServiceEntity,
) CreatedEvent {
	return CreatedEvent{
		Record:   appointment,
		Customer: customer,
		Service:  service,
	}
}

func (e CreatedEvent) Type() EventType {
	return CreatedEventType
}

type CanceledEvent struct {
	Record   RecordEntity
	Customer CustomerEntity
	Service  ServiceEntity
}

func NewAppointmentCanceled(
	appointment RecordEntity,
	customer CustomerEntity,
	service ServiceEntity,
) CanceledEvent {
	return CanceledEvent{
		Record:   appointment,
		Customer: customer,
		Service:  service,
	}
}

func (e CanceledEvent) Type() EventType {
	return CanceledEventType
}

type ChangeType int

const (
	CreatedChangeType ChangeType = iota
	StatusChangeType
	DateTimeChangeType
	RemovedChangeType
)

type ChangedEvent struct {
	ChangeType ChangeType
	Record     RecordEntity
}

func NewChanged(
	changeType ChangeType,
	appointment RecordEntity,
) ChangedEvent {
	return ChangedEvent{
		ChangeType: changeType,
		Record:     appointment,
	}
}

func (e ChangedEvent) Type() EventType {
	return ChangedEventType
}
