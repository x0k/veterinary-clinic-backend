//go:build js && wasm

package main

import (
	_ "github.com/x0k/veterinary-clinic-backend/internal/appointment/module/wasm"
)

func main() {}
