//go:build js && wasm

package js_adapters

import (
	"context"
	"errors"
	"syscall/js"
)

var PromiseConstructor = js.Global().Get("Promise")
var ErrorConstructor = js.Global().Get("Error")
var ObjectConstructor = js.Global().Get("Object")
var Console = js.Global().Get("console")

func NewPromise(action func() (js.Value, *js.Value)) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]
		go func() {
			data, err := action()
			if err != nil {
				reject.Invoke(*err)
			} else {
				resolve.Invoke(data)
			}
		}()
		return nil
	})
	promise := PromiseConstructor.New(handler)
	handler.Release()
	return promise
}

func Promise(action func() (js.Value, error)) js.Value {
	return NewPromise(func() (js.Value, *js.Value) {
		res, err := action()
		if err != nil {
			jsErr := NewError(err)
			return js.Undefined(), &jsErr
		}
		return res, nil
	})
}

func Resolve(data js.Value) js.Value {
	return PromiseConstructor.Invoke("resolve", data)
}

func Reject(err js.Value) js.Value {
	return PromiseConstructor.Invoke("reject", err)
}

func NewError(err error) js.Value {
	return ErrorConstructor.New(err.Error())
}

func RejectError(err error) js.Value {
	return Reject(NewError(err))
}

func Await(ctx context.Context, promise js.Value) (js.Value, error) {
	resChan := make(chan js.Value)
	errChan := make(chan error)
	go func() {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
			onSuccess := js.FuncOf(func(this js.Value, args []js.Value) any {
				resChan <- args[0]
				return nil
			})
			onError := js.FuncOf(func(this js.Value, args []js.Value) any {
				errChan <- errors.New(args[0].Call("toString").String())
				return nil
			})
			var finally js.Func
			finally = js.FuncOf(func(this js.Value, args []js.Value) any {
				onSuccess.Release()
				onError.Release()
				finally.Release()
				return nil
			})
			promise.Call("then", onSuccess, onError).Call("finally", finally)
		}
	}()
	select {
	case <-ctx.Done():
		return js.Undefined(), ctx.Err()
	case err := <-errChan:
		return js.Undefined(), err
	case data := <-resChan:
		return data, nil
	}
}
