package main

import (
	"fmt"
	"syscall/js"
)

func registerCallbacks(ws js.Value) {
	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		message := args[0].Get("data").String()
		fmt.Println("message registerCallbacks ")

		fmt.Println(message)
		return nil
	}))
}

func sendCallbacks(ws js.Value) {
	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		ws.Call("send", js.ValueOf("{cmd}"))
		/*		message := args[0].Get("data").String()
				fmt.Println("message sendCallbacks ")
				fmt.Println(message)*/
		return nil
	}))
}

func main() {
	c := make(chan struct{}, 0)
	ws := js.Global().Get("WebSocket").New("ws://127.0.0.1:8080")
	sendCallbacks(ws)
	registerCallbacks(ws)
	<-c
}
