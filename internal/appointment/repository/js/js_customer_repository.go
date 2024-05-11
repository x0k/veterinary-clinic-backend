//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type CustomerRepositoryConfig struct {
	CreateCustomer     *js.Value `js:"createCustomer"`
	CustomerByIdentity *js.Value `js:"loadCustomerByIdentity"`
	UpdateCustomer     *js.Value `js:"updateCustomer"`
}

type CustomerRepository struct {
	cfg CustomerRepositoryConfig
}

func NewCustomerRepository(
	cfg CustomerRepositoryConfig,
) *CustomerRepository {
	return &CustomerRepository{
		cfg: cfg,
	}
}

func (r *CustomerRepository) CreateCustomer(
	ctx context.Context,
	customer *appointment.CustomerEntity,
) error {
	promise := r.cfg.CreateCustomer.Invoke(
		vert.ValueOf(appointment_js_adapters.CustomerToDTO(*customer)),
	)
	customerId, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return err
	}
	customer.SetId(appointment.NewCustomerId(customerId.String()))
	return nil
}

func (r *CustomerRepository) CustomerByIdentity(
	ctx context.Context,
	identity appointment.CustomerIdentity,
) (appointment.CustomerEntity, error) {
	promise := r.cfg.CustomerByIdentity.Invoke(
		identity.String(),
	)
	jsValue, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return appointment.CustomerEntity{}, err
	}
	if jsValue.IsNull() {
		return appointment.CustomerEntity{}, shared.ErrNotFound
	}
	var dto appointment_js_adapters.CustomerDTO
	if err := vert.Assign(jsValue, &dto); err != nil {
		return appointment.CustomerEntity{}, err
	}
	return appointment_js_adapters.CustomerFromDTO(dto)
}

func (r *CustomerRepository) UpdateCustomer(
	ctx context.Context,
	customer appointment.CustomerEntity,
) error {
	promise := r.cfg.UpdateCustomer.Invoke(
		vert.ValueOf(appointment_js_adapters.CustomerToDTO(customer)),
	)
	_, err := js_adapters.Await(ctx, promise)
	return err
}
