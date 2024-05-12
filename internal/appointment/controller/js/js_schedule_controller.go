//go:build js && wasm

package appointment_js_controller

import (
	"context"
	"syscall/js"
	"time"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
	appointment_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case"
	appointment_js_use_case "github.com/x0k/veterinary-clinic-backend/internal/appointment/use_case/js"
)

func NewSchedule(
	ctx context.Context,
	module js.Value,
	scheduleUseCase *appointment_use_case.ScheduleUseCase[js_adapters.Result],
	dayOrNextWorkingDayUseCase *appointment_js_use_case.DayOrNextWorkingDayUseCase[js_adapters.Result],
	upsertCustomerUseCase *appointment_js_use_case.UpsertCustomerUseCase[js_adapters.Result],
) {
	module.Set("schedule", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		preferredDate := args[0].String()
		date, err := time.Parse(time.RFC3339, preferredDate)
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return scheduleUseCase.Schedule(ctx, time.Now(), date)
		})
	}))
	module.Set("dayOrNextWorkingDay", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		now := args[0].String()
		date, err := time.Parse(time.RFC3339, now)
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return dayOrNextWorkingDayUseCase.DayOrNextWorkingDay(ctx, date)
		})
	}))
	module.Set("upsertCustomer", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		var createCustomerDTO appointment_js_adapters.CreateCustomerDTO
		if err := vert.Assign(args[0], &createCustomerDTO); err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return upsertCustomerUseCase.Upsert(
				ctx,
				createCustomerDTO.IdentityProvider,
				createCustomerDTO.Identity,
				createCustomerDTO.Name,
				createCustomerDTO.Phone,
				createCustomerDTO.Email,
			)
		})
	}))
}
