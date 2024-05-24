//go:build js && wasm

package appointment_js_presenter

import (
	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

func ServicesPresenter(services []appointment.ServiceEntity) (js_adapters.Result, error) {
	servicesDto := make([]appointment_js_adapters.ServiceDTO, len(services))
	var err error
	for i, service := range services {
		servicesDto[i], err = appointment_js_adapters.ServiceToDTO(service)
		if err != nil {
			return js_adapters.Result{}, err
		}
	}
	return js_adapters.Ok(vert.ValueOf(servicesDto)), nil
}
