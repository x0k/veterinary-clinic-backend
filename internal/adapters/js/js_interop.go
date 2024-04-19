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

func Promise(action func() (js.Value, *js.Value)) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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

func Resolve(data js.Value) js.Value {
	return PromiseConstructor.Invoke("resolve", data)
}

func Reject(err js.Value) js.Value {
	return PromiseConstructor.Invoke("reject", err)
}

func Error(err error) js.Value {
	return ErrorConstructor.New(err.Error())
}

func RejectError(err error) js.Value {
	return Reject(Error(err))
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
			// Wait for the promise to resolve
			onSuccess := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resChan <- args[0]
				return nil
			})
			onError := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				errChan <- errors.New(args[0].Invoke("toString").String())
				return nil
			})
			promise.Call("then", onSuccess, onError)
			onSuccess.Release()
			onError.Release()
		}
	}()

	select {
	case <-ctx.Done():
		return js.Null(), ctx.Err()
	case err := <-errChan:
		return js.Null(), err
	case data := <-resChan:
		return data, nil
	}
}
