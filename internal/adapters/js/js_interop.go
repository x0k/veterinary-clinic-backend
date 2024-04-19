//go:build js && wasm

package js_adapters

import "syscall/js"

var PromiseConstructor = js.Global().Get("Promise")
var ErrorConstructor = js.Global().Get("Error")
var ObjectConstructor = js.Global().Get("Object")

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

	return PromiseConstructor.New(handler)
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
