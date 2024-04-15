package appointment

import "context"

type NotificationSender[R any] interface {
	SendNotification(ctx context.Context, message R) error
}
