package controller

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

var (
	RegisterBody = &Message{CMD: 0}
	CloseBody    = &Message{CMD: 1}
)

func ReceiveCallbacks(ws js.Value) {
	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		body := args[0].Get("data").String()
		fmt.Println("message ReceiveCallbacks: ")
		fmt.Println(body)
		var message = &Message{}
		json.Unmarshal([]byte(body), message)
		messageHandler(ws, message)
		return nil
	}))
}

func RegisterCallbacks(ws js.Value) {
	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.Call("send", js.ValueOf(RegisterBody.String()))
		return nil
	}))
}

func ErrorCallbacks(ws js.Value) {
	ws.Call("addEventListener", "error", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		alert := js.Global().Get("alert")
		alert.Invoke("error")
		return nil
	}))
}

func CloseCallbacks(ws js.Value) {
	ws.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		alert := js.Global().Get("alert")
		alert.Invoke("close")
		return nil
	}))
}

func messageHandler(ws js.Value, message *Message) {
	switch message.CMD {
	case 3:
		//todo something
		fmt.Println("Assign task:" + message.Payload.TaskId)
	case 4:
		//todo something
		fmt.Println("Interrupt task:" + message.Payload.TaskId)
	}
}
