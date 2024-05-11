package appointment_js_adapters

import "github.com/x0k/veterinary-clinic-backend/internal/appointment"

type CustomerIdentityProvider string

const (
	VkIdentityProvider CustomerIdentityProvider = "vk"
)

type CreateCustomerDTO struct {
	IdentityProvider CustomerIdentityProvider `js:"identityProvider"`
	Identity         string                   `js:"identity"`
	Name             string                   `js:"name"`
	Phone            string                   `js:"phone"`
	Email            string                   `js:"email"`
}

type CustomerDTO struct {
	Id       string `js:"id"`
	Identity string `js:"identity"`
	Name     string `js:"name"`
	Phone    string `js:"phone"`
	Email    string `js:"email"`
}

func CustomerToDTO(customer appointment.CustomerEntity) CustomerDTO {
	return CustomerDTO{
		Id:       customer.Id.String(),
		Identity: customer.Identity.String(),
		Name:     customer.Name,
		Phone:    customer.PhoneNumber,
		Email:    customer.Email,
	}
}

func CustomerFromDTO(dto CustomerDTO) (appointment.CustomerEntity, error) {
	identity, err := appointment.NewCustomerIdentity(dto.Identity)
	if err != nil {
		return appointment.CustomerEntity{}, err
	}
	return appointment.CustomerEntity{
		Id:          appointment.NewCustomerId(dto.Id),
		Identity:    identity,
		Name:        dto.Name,
		PhoneNumber: dto.Phone,
		Email:       dto.Email,
	}, nil
}
