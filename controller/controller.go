package controller

import (
	"fmt"
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
	"google.golang.org/protobuf/proto"
	"syscall/js"
	"time"
	"worker/executor"
)

var (
	RegisterBody = &Msg{Cmd: CMD_Register}
	CloseBody    = &Msg{Cmd: CMD_Close}
)

func ReceiveCallbacks(ws js.Value) {
	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		buffer := args[0]
		uint8Array := js.Global().Get("Uint8Array").New(buffer)
		data := make([]byte, uint8Array.Get("length").Int())
		js.CopyBytesToGo(data, uint8Array)
		fmt.Println("message ReceiveCallbacks: ")
		var receiveMsg = &Msg{}
		err := proto.Unmarshal(data, receiveMsg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("cmd is: %v\n", receiveMsg.Cmd)
		//TODO func builder gen hash Coins
		messageHandler(ws, receiveMsg, executor.NewHashCoins(receiveMsg))
		return nil
	}))
}

func sendMsg(ws js.Value, msg *Msg) {
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	dest := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(dest, data)
	ws.Call("send", dest)
}

func RegisterCallbacks(ws js.Value) {
	ws.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("register worker")
		sendMsg(ws, RegisterBody)
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
		ws.Call("send", js.ValueOf(CloseBody.String()))
		alert := js.Global().Get("alert")
		alert.Invoke("close")
		return nil
	}))
}

func messageHandler(ws js.Value, message *Msg, exec executor.Executor) {
	switch message.Cmd {
	case CMD_Register:
		fmt.Println("Registered successfully")
	case CMD_Assign:
		exec.Start()
		fmt.Println("Assign task:" + message.GetAssign().TaskId)
	case CMD_Close:
		msg := exec.Interrupt()
		sendMsg(ws, msg)
		fmt.Println("Interrupt task:" + message.GetAssign().TaskId)
	default:
		// CMD_Status
		for {
			time.Sleep(time.Second * 2)
			msg := exec.Status()
			sendMsg(ws, msg)
		}
	}

}
