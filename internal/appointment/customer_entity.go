package appointment

import (
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type CustomerIdentity string

func NewTelegramCustomerIdentity(id entity.TelegramUserId) CustomerIdentity {
	return CustomerIdentity(fmt.Sprintf("tg-%d", id))
}

type CustomerId string

func NewCustomerId(id string) CustomerId {
	return CustomerId(id)
}

func (c CustomerId) String() string {
	return string(c)
}

type CustomerEntity struct {
	Id          CustomerId
	Name        string
	PhoneNumber string
	Email       string
}

func NewCustomer(
	id CustomerId,
	name string,
	phoneNumber string,
	email string,
) CustomerEntity {
	return CustomerEntity{
		Id:          id,
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}
