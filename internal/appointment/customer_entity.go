package appointment

import (
	"fmt"
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrUnknownCustomerIdType = fmt.Errorf("unknown customer id type")

type CustomerIdType string

const (
	TelegramIdType CustomerIdType = "tg"
	VkIdType       CustomerIdType = "vk"
)

type CustomerId string

func NewCustomerId(id string) CustomerId {
	return CustomerId(id)
}

func NewTelegramCustomerId(id entity.TelegramUserId) CustomerId {
	return CustomerId(fmt.Sprintf("%s-%d", TelegramIdType, id))
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

func (c *CustomerEntity) IdType() (CustomerIdType, error) {
	if strings.HasPrefix(c.Id.String(), string(TelegramIdType)) {
		return TelegramIdType, nil
	}
	if strings.HasPrefix(c.Id.String(), string(VkIdType)) {
		return VkIdType, nil
	}
	return "", ErrUnknownCustomerIdType
}

func TelegramUserIdToCustomerId(id entity.TelegramUserId) CustomerId {
	return NewCustomerId(fmt.Sprintf("tg-%d", id))
}
