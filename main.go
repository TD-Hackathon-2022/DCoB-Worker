package main

import (
	"syscall/js"
	"worker/controller"
)

func main() {
	ws := js.Global().Get("WebSocket").New("ws://127.0.0.1:8080")
	controller.RegisterCallbacks(ws)
	controller.ReceiveCallbacks(ws)
	controller.CloseCallbacks(ws)
	controller.ErrorCallbacks(ws)
	select {}
}
