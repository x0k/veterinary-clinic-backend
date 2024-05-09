//go:build js && wasm

package js_adapters

import (
	"syscall/js"
)

var ErrorConstructor = js.Global().Get("Error")

func Error(err error) js.Value {
	return ErrorConstructor.New(err.Error())
}
