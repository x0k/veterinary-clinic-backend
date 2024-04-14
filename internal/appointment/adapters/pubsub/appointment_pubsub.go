package appointment_pubsub_adapters

import (
	pubsub_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/pubsub"
	appointment_event "github.com/x0k/veterinary-clinic-backend/internal/appointment/event"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

func Subscribe[E pubsub.Event[appointment_event.Type]](
	subs pubsub.SubscriptionsManager[appointment_event.Type],
	preStopper module.PreStopper,
) <-chan E {
	return pubsub_adapters.Subscribe[appointment_event.Type, E](subs, preStopper)
}
