package entity

import (
	"fmt"
)

type UserId string

type User struct {
	Id          UserId
	Name        string
	PhoneNumber string
	Email       string
}

type TelegramUserId int64
type TelegramUsername string
type TelegramFirstName string
type TelegramLastName string

func TelegramUserIdToUserId(id TelegramUserId) UserId {
	return UserId(fmt.Sprintf("tg-%d", id))
}

func NewTelegramUser(
	id TelegramUserId,
	username TelegramUsername,
	firstName TelegramFirstName,
	lastName TelegramLastName,
) User {
	return User{
		Id:          TelegramUserIdToUserId(id),
		Name:        fmt.Sprintf("%s %s", firstName, lastName),
		PhoneNumber: "",
		Email:       fmt.Sprintf("@%s", username),
	}
}
