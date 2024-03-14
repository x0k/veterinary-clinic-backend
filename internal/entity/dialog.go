package entity

type DialogId string

type Dialog struct {
	Id     DialogId
	UserId UserId
}

type DialogMessage[R any] struct {
	DialogId DialogId
	Message  R
}
