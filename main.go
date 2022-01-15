package worker

import (
	"fmt"
	"syscall/js"
)

func test() {

	ws := js.Global().Get("WebSocket").New("ws://localhost:8080/ws")

	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("open")
		ws.Call("send", js.ValueOf([]byte{123}))
		return nil
	}))
}

func main() {
	test()
}
