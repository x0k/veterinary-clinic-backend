//go:build js && wasm

package appointment_js_controller

import (
	"context"
	"syscall/js"
	"time"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_js_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/js"
)

func NewSchedule(
	module js.Value,
	scheduleUseCase *appointment_use_case.ScheduleUseCase[js_adapters.Result],
	dayOrNextWorkingDayUseCase *appointment_js_use_case.DayOrNextWorkingDayUseCase[js_adapters.Result],
) {
	module.Set("schedule", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		ctx := context.TODO()
		preferredDate := args[0].String()
		date, err := time.Parse(time.RFC3339, preferredDate)
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return scheduleUseCase.Schedule(ctx, time.Now(), date)
		})
	}))
}
