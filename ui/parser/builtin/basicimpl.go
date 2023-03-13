//go:build js

package builtin

import (
	"log"
	"syscall/js"
)

func ToggleSingle(this js.Value, arg []js.Value) any {
	log.Printf("xxx got a call to toggle single: %v", this)
	return nil
}
