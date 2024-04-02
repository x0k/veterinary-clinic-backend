package appointment

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ClientIdentity string

func NewTelegramClientIdentity(id entity.TelegramUserId) ClientIdentity {
	return ClientIdentity(fmt.Sprintf("tg-%d", id))
}

type ClientId string

func (c ClientId) String() string {
	return string(c)
}

type ClientEntity struct {
	Id          ClientId
	Name        string
	PhoneNumber string
	Email       string
}

func NewClient(
	name string,
	phoneNumber string,
	email string,
) *ClientEntity {
	return &ClientEntity{
		Id:          ClientId(uuid.New().String()),
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}
