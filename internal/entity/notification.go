package entity

type NotificationMessage[R any] struct {
	UserId  UserId
	Message R
}
