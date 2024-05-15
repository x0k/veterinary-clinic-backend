package appointment_js_use_case

import (
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrInvalidIdentityProvider = errors.New("invalid identity provider")

func customerIdentity(
	userIdentityProvider appointment_js_adapters.CustomerIdentityProvider,
	userIdentity string,
) (appointment.CustomerIdentity, error) {
	switch userIdentityProvider {
	case appointment_js_adapters.VkIdentityProvider:
		return appointment.NewVkCustomerIdentity(
			shared.NewVkUserId(userIdentity),
		)
	default:
		return "", ErrInvalidIdentityProvider
	}
}
