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
	module.Set("timePeriodDurationInMinutes", js_adapters.Sync(func(args []js.Value) js_adapters.Result {
		if len(args) < 1 {
			return js_adapters.Fail(js_adapters.ErrTooFewArguments)
		}
		var timePeriodDto shared_js_adapters.TimePeriodDTO
		if err := vert.Assign(args[0], &timePeriodDto); err != nil {
			return js_adapters.Fail(err)
		}
		return js_adapters.Ok(
			js.ValueOf(
				shared.TimePeriodDurationInMinutes(
					shared_js_adapters.TimePeriodFromDTO(timePeriodDto),
				).Int(),
			),
		)
	}))
	module.Set("isDateTimePeriodIntersectWithPeriods", js_adapters.Sync(func(args []js.Value) js_adapters.Result {
		if len(args) < 2 {
			return js_adapters.Fail(js_adapters.ErrTooFewArguments)
		}
		var periodDto shared_js_adapters.DateTimePeriodDTO
		if err := vert.Assign(args[0], &periodDto); err != nil {
			return js_adapters.Fail(err)
		}
		period := shared_js_adapters.DateTimePeriodFromDTO(periodDto)
		var el shared_js_adapters.DateTimePeriodDTO
		for i := 1; i < len(args); i++ {
			if err := vert.Assign(args[i], &el); err != nil {
				return js_adapters.Fail(err)
			}
			if shared.DateTimePeriodApi.IsValidPeriod(
				shared.DateTimePeriodApi.IntersectPeriods(
					period,
					shared_js_adapters.DateTimePeriodFromDTO(el),
				),
			) {
				return js_adapters.Ok(
					js.ValueOf(
						true,
					),
				)
			}

		}
		return js_adapters.Ok(
			js.ValueOf(
				false,
			),
		)
	}))
}
