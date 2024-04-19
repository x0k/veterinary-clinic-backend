//go:build js && wasm

package js_adapters

type ConsoleLoggerHandler struct {
	disabled bool
}

func (h *ConsoleLoggerHandler) Enabled() bool {
	return !h.disabled
}
