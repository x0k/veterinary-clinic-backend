//go:build js && wasm

package shared_js_controller

import (
	"syscall/js"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	shared_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/shared/adapters/js"
)

func NewDateTime(
	module js.Value,
) {
	module.Set("timePeriodDurationInMinutes", js.FuncOf(func(this js.Value, args []js.Value) any {
		var timePeriodDto shared_js_adapters.TimePeriodDTO
		if err := vert.Assign(args[0], &timePeriodDto); err != nil {
			return js_adapters.RejectError(err)
		}
		return shared.TimePeriodDurationInMinutes(shared_js_adapters.TimePeriodFromDTO(timePeriodDto)).Minutes()
	}))
}
