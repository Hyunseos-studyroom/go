package main

import (
	"syscall/js"
)

func add(this js.Value, p []js.Value) interface{} {
	return js.ValueOf(p[0].Int() + p[1].Int())
}

func registerCallbacks() {
	js.Global().Set("add", js.FuncOf(add))
}

func main() {
	c := make(chan struct{}, 0)

	registerCallbacks()

	<-c // Prevent the program from exiting
}
