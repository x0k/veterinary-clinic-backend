package appointment_event

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type appointmentCreatedHandler chan<- AppointmentCreatedEvent

func (h appointmentCreatedHandler) Type() Type {
	return AppointmentCreated
}

func (h appointmentCreatedHandler) Handle(event pubsub.Event[Type]) {
	h <- event.(AppointmentCreatedEvent)
}

type Publisher struct {
	pubSub *pubsub.PubSub[Type]
}

func NewPublisher() *Publisher {
	return &Publisher{
		pubSub: pubsub.New[Type](),
	}
}

func (p *Publisher) HandleAppointmentCreate(c chan<- AppointmentCreatedEvent) func() {
	return func() {}
}
