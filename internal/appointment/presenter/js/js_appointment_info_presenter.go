//go:build js && wasm

package appointment_js_presenter

import (
	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

func AppointmentInfoPresenter(
	app appointment.RecordEntity,
	service appointment.ServiceEntity,
) (js_adapters.Result, error) {
	return js_adapters.Ok(
		vert.ValueOf(
			appointment_js_adapters.AppointmentInfoToDTO(
				app,
				service,
			),
		),
	), nil
}
