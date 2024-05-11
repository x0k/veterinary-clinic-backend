package appointment

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrUnknownCustomerIdentityType = errors.New("unknown customer identity type")
var ErrWrongCustomerIdentityType = errors.New("wrong customer identity type")
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

func NewCustomerIdentity(identity string) (CustomerIdentity, error) {
	if !strings.HasPrefix(identity, TelegramIdentityType.String()) ||
		!strings.HasPrefix(identity, VkIdentityType.String()) {
		return "", ErrUnknownCustomerIdentityType
	}
	return CustomerIdentity(identity), nil
}

func NewTelegramCustomerIdentity(id shared.TelegramUserId) (CustomerIdentity, error) {
	return NewCustomerIdentity(fmt.Sprintf("%s-%d", TelegramIdentityType, id))
}

func NewVkCustomerIdentity(id shared.VkUserId) (CustomerIdentity, error) {
	return NewCustomerIdentity(fmt.Sprintf("%s-%s", VkIdentityType, id))
}

func (identity CustomerIdentity) ToTelegramUserId() (shared.TelegramUserId, error) {
	tp, err := identity.Type()
	if err != nil {
		return 0, err
	}
	if tp != TelegramIdentityType {
		return 0, ErrWrongCustomerIdentityType
	}
	idStr := strings.Split(string(identity), "-")[1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return shared.NewTelegramUserId(id), nil
}

func (identity CustomerIdentity) String() string {
	return string(identity)
}

func (identity CustomerIdentity) Type() (CustomerIdentityType, error) {
	if strings.HasPrefix(identity.String(), TelegramIdentityType.String()) {
		return TelegramIdentityType, nil
	}
	if strings.HasPrefix(identity.String(), VkIdentityType.String()) {
		return VkIdentityType, nil
	}
	return "", ErrUnknownCustomerIdentityType
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
	return c.Identity.Type()
}

func (c *CustomerEntity) Update(
	name string,
	phoneNumber string,
	email string,
) bool {
	updated := false
	if c.Name != name {
		c.Name = name
		updated = true
	}
	if c.PhoneNumber != phoneNumber {
		c.PhoneNumber = phoneNumber
		updated = true
	}
	if c.Email != email {
		c.Email = email
		updated = true
	}
	return updated
}
