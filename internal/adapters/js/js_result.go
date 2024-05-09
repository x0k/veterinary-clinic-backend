//go:build js && wasm

package js_adapters

import "syscall/js"

func NewOk(value any) js.Value {
	obj := ObjectConstructor.New()
	obj.Set("ok", true)
	obj.Set("value", value)
	return obj
}

func Ok(value js.Value) js.Value {
	return NewOk(value)
}

func NewFail(err any) js.Value {
	obj := ObjectConstructor.New()
	obj.Set("ok", false)
	obj.Set("error", err)
	return obj
}

func Fail(err error) js.Value {
	return NewFail(NewError(err))
}
