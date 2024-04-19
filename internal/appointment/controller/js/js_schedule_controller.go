//go:build js && wasm

package appointment_js_controller

import (
	"context"
	"syscall/js"
	"time"

	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
)

func NewSchedule(
	module js.Value,
	scheduleUseCase *appointment_use_case.ScheduleUseCase[js.Value],
) {
	module.Set("schedule", js.FuncOf(func(this js.Value, args []js.Value) any {
		ctx := context.TODO()
		preferredDate := args[0].String()
		date, err := time.Parse(time.RFC3339, preferredDate)
		if err != nil {
			return js_adapters.RejectError(err)
		}
		schedule, err := scheduleUseCase.Schedule(ctx, time.Now(), date)
		if err != nil {
			return js_adapters.RejectError(err)
		}
		return schedule
	}))
}
