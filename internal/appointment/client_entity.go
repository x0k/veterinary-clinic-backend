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

type ClientId uuid.UUID

func (c ClientId) String() string {
	return uuid.UUID(c).String()
}

type Client struct {
	Id          ClientId
	Name        string
	PhoneNumber string
	Email       string
}

func NewClient(
	name string,
	phoneNumber string,
	email string,
) *Client {
	return &Client{
		Id:          ClientId(uuid.New()),
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}
