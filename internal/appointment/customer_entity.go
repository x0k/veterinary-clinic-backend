package appointment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrUnknownCustomerIdentityType = errors.New("unknown customer identity type")
var ErrInvalidCustomerId = errors.New("invalid customer id")
var ErrCustomerIdIsNotTemporal = errors.New("id is not temporal")

type CustomerIdentityType string

func (t CustomerIdentityType) String() string {
	return string(t)
}

const (
	TelegramIdentityType CustomerIdentityType = "tg"
	VkIdentityType       CustomerIdentityType = "vk"
)

type CustomerIdentity string

func NewCustomerIdentity(identity string) CustomerIdentity {
	return CustomerIdentity(identity)
}

func NewTelegramCustomerIdentity(id entity.TelegramUserId) CustomerIdentity {
	return CustomerIdentity(fmt.Sprintf("%s-%d", TelegramIdentityType, id))
}

func (identity CustomerIdentity) String() string {
	return string(identity)
}

type CustomerId string

const TemporalCustomerId CustomerId = "tmp_customer_id"

func NewCustomerId(id string) CustomerId {
	return CustomerId(id)
}

func (id CustomerId) String() string {
	return string(id)
}

type CustomerEntity struct {
	Id          CustomerId
	Identity    CustomerIdentity
	Name        string
	PhoneNumber string
	Email       string
}

func NewCustomer(
	id CustomerId,
	identity CustomerIdentity,
	name string,
	phoneNumber string,
	email string,
) CustomerEntity {
	return CustomerEntity{
		Id:          id,
		Identity:    identity,
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}

func (c *CustomerEntity) SetId(id CustomerId) error {
	if c.Id != TemporalCustomerId {
		return fmt.Errorf("%w: %s", ErrCustomerIdIsNotTemporal, c.Id)
	}
	c.Id = id
	return nil
}

func (c *CustomerEntity) IdentityType() (CustomerIdentityType, error) {
	if strings.HasPrefix(c.Identity.String(), TelegramIdentityType.String()) {
		return TelegramIdentityType, nil
	}
	if strings.HasPrefix(c.Identity.String(), VkIdentityType.String()) {
		return VkIdentityType, nil
	}
	return "", ErrUnknownCustomerIdentityType
}
