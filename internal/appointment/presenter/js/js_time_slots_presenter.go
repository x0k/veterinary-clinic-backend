//go:build js && wasm

package appointment_js_presenter

import (
	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	shared_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/shared/adapters/js"
)

func FreeTimeSlotsPresenter(
	slots appointment.SampledFreeTimeSlots,
) (js_adapters.Result, error) {
	periods := make([]shared_js_adapters.TimePeriodDTO, len(slots))
	for i, s := range slots {
		periods[i] = shared_js_adapters.TimePeriodToDTO(s)
	}
	return js_adapters.Ok(vert.ValueOf(periods)), nil
}
