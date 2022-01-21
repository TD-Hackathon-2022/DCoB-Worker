package controller

import (
	"fmt"
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
	"google.golang.org/protobuf/proto"
	"sync"
	"syscall/js"
	. "worker/executor"
	"worker/htmlprinter"
)

var (
	RegisterBody = &Msg{Cmd: CMD_Register, Payload: &Msg_Empty{Empty: &EmptyPayload{}}}
	CloseBody    = &Msg{Cmd: CMD_Close, Payload: &Msg_Empty{Empty: &EmptyPayload{}}}
)

func ReceiveCallbacks(ws js.Value) {
	ws.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Get("data").Call("arrayBuffer").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			uint8Array := js.Global().Get("Uint8Array").New(args[0])
			data := make([]byte, uint8Array.Get("length").Int())
			js.CopyBytesToGo(data, uint8Array)
			var receiveMsg = &Msg{}
			err := proto.Unmarshal(data, receiveMsg)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("ReceiveCallbacks recieved: %v\n", receiveMsg)
			//	htmlprinter.PrintPHtml(receiveMsg.String())
			messageHandler(ws, receiveMsg)
			return nil
		}))
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
		alert.Invoke("server error!!")
		return nil
	}))
}

func CloseCallbacks(ws js.Value) {
	ws.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.Call("send", js.ValueOf(CloseBody.String()))
		alert := js.Global().Get("alert")
		alert.Invoke("server closed!!")
		return nil
	}))
}

var exec Executor = nil
var lock sync.RWMutex

func messageHandler(ws js.Value, message *Msg) {
	switch message.Cmd {
	case CMD_Register:
		fmt.Println("Registered successfully")
		htmlprinter.PrintHHtml("worker registered successfully!!")
	case CMD_Interrupt:
		fmt.Printf("Interrupt task: %s\n", message.GetInterrupt().TaskId)
		lock.RLock()
		lock.RUnlock()
		if exec != nil {
			msg := exec.Interrupt()
			sendMsg(ws, msg)
		}
		htmlprinter.PrintHHtml("interrupt task:" + message.GetInterrupt().TaskId)
	case CMD_Assign:
		fmt.Printf("Assign task: %s\n", message.GetAssign().TaskId)
		e := ExecutorBuilder(message)
		lock.Lock()
		defer lock.Unlock()
		if exec != nil && exec.Status().GetStatus().TaskStatus == TaskStatus_Running {
			result := fmt.Sprintf("Cannot assign task: %s, current task %s still running!",
				message.GetAssign().TaskId,
				e.Status().GetStatus().TaskId)
			fmt.Println(result)
			sendMsg(ws, &Msg{
				Cmd: CMD_Status,
				Payload: &Msg_Status{
					Status: &StatusPayload{
						WorkStatus: WorkerStatus_Busy,
						TaskId:     message.GetAssign().TaskId,
						TaskStatus: TaskStatus_Error,
						ExecResult: result,
					},
				},
			})
			return
		}

		exec = e
		e.Start(func() {
			fmt.Printf("Task finished: %s\n", message.GetAssign().TaskId)
			sendMsg(ws, e.Status())
			htmlprinter.PrintHHtml("Assign task: " + message.GetAssign().TaskId)
		})
	}
}
