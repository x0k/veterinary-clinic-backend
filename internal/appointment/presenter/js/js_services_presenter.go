//go:build js && wasm

package appointment_js_presenter

import (
	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

func ServicesPresenter(services []appointment.ServiceEntity) (js_adapters.Result, error) {
	dtos := make([]appointment_js_adapters.ServiceDTO, len(services))
	for i, service := range services {
		dtos[i] = appointment_js_adapters.ServiceToDTO(service)
	}
	return js_adapters.Ok(vert.ValueOf(dtos)), nil
}
