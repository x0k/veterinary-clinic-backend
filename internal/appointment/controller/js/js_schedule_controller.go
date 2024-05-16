//go:build js && wasm

package appointment_js_controller

import (
	"context"
	"syscall/js"
	"time"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
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
	freeTimeSlotsUseCase *appointment_js_use_case.FreeTimeSlotsUseCase[js_adapters.Result],
	activeAppointmentUseCase *appointment_js_use_case.ActiveAppointmentUseCase[js_adapters.Result],
	createAppointmentUseCase *appointment_use_case.MakeAppointmentUseCase[js_adapters.Result],
) {
	module.Set("schedule", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		if len(args) < 1 {
			return js_adapters.ResolveError(js_adapters.ErrTooFewArguments)
		}
		date, err := time.Parse(time.RFC3339, args[0].String())
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return scheduleUseCase.Schedule(ctx, time.Now(), date)
		})
	}))
	module.Set("dayOrNextWorkingDay", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		if len(args) < 1 {
			return js_adapters.ResolveError(js_adapters.ErrTooFewArguments)
		}
		date, err := time.Parse(time.RFC3339, args[0].String())
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return dayOrNextWorkingDayUseCase.DayOrNextWorkingDay(ctx, date)
		})
	}))
	module.Set("upsertCustomer", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		if len(args) < 1 {
			return js_adapters.ResolveError(js_adapters.ErrTooFewArguments)
		}
		var createCustomerDTO appointment_js_adapters.CreateCustomerDTO
		if err := vert.Assign(args[0], &createCustomerDTO); err != nil {
			return js_adapters.ResolveError(err)
		}
		identity, err := appointment.NewCustomerIdentity(createCustomerDTO.Identity)
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return upsertCustomerUseCase.Upsert(
				ctx,
				identity,
				createCustomerDTO.Name,
				createCustomerDTO.Phone,
				createCustomerDTO.Email,
			)
		})
	}))
	module.Set("freeTimeSlots", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		if len(args) < 2 {
			return js_adapters.ResolveError(js_adapters.ErrTooFewArguments)
		}
		serviceId := appointment.NewServiceId(args[0].String())
		appointmentDate, err := time.Parse(time.RFC3339, args[1].String())
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return freeTimeSlotsUseCase.FreeTimeSlots(
				ctx,
				serviceId,
				time.Now(),
				appointmentDate,
			)
		})
	}))
	module.Set("activeAppointment", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		if len(args) < 1 {
			return js_adapters.ResolveError(js_adapters.ErrTooFewArguments)
		}
		identity, err := appointment.NewCustomerIdentity(args[0].String())
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return activeAppointmentUseCase.ActiveAppointment(
				ctx,
				identity,
			)
		})
	}))
	module.Set("createAppointment", js_adapters.Async(func(args []js.Value) js_adapters.Promise {
		if len(args) < 3 {
			return js_adapters.ResolveError(js_adapters.ErrTooFewArguments)
		}
		appointmentDate, err := time.Parse(time.RFC3339, args[0].String())
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		customerIdentity, err := appointment.NewCustomerIdentity(args[1].String())
		if err != nil {
			return js_adapters.ResolveError(err)
		}
		serviceId := appointment.NewServiceId(args[2].String())
		return js_adapters.NewPromise(func() (js_adapters.Result, error) {
			return createAppointmentUseCase.CreateAppointment(
				ctx,
				time.Now(),
				appointmentDate,
				customerIdentity,
				serviceId,
			)
		})
	}))
}
