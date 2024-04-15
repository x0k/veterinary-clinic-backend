package appointment_pubsub_controller

import (
	pubsub_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/pubsub"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

func Subscribe[E pubsub.Event[appointment.EventType]](
	subs pubsub.SubscriptionsManager[appointment.EventType],
	preStopper module.PreStopper,
) <-chan E {
	return pubsub_adapters.Subscribe[appointment.EventType, E](subs, preStopper)
}
