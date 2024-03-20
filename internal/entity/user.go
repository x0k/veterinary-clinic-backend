package entity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrNotTelegramUser = errors.New("not a telegram user")

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

func IsTelegramUserId(userId UserId) bool {
	return strings.HasPrefix("tg-", string(userId))
}

func UserIdToTelegramUserId(userId UserId) (TelegramUserId, error) {
	if !IsTelegramUserId(userId) {
		return 0, ErrNotTelegramUser
	}
	id, err := strconv.ParseInt(string(userId)[3:], 10, 64)
	if err != nil {
		return 0, err
	}
	return TelegramUserId(id), nil
}
