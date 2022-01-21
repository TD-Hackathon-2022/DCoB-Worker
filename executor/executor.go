package executor

import (
	"encoding/base64"
	"fmt"
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
	"strings"
	"syscall/js"
)

type Executor interface {
	Start(thenDo func())
	Status() *Msg
	Interrupt() *Msg
}

const (
	hashMiner = "hash-miner"
	custom    = "custom-func-"
)

func ExecutorBuilder(msg *Msg) Executor {
	funcId := msg.GetAssign().FuncId
	if strings.HasPrefix(funcId, custom) {
		funcPrefixArr := strings.Split(funcId, custom)
		content, _ := base64.StdEncoding.DecodeString(msg.GetAssign().Data)
		return newCustomFunc(funcPrefixArr[1], msg, content)
	}

	switch funcId {
	case hashMiner:
		return NewHashCoins(msg)
	default:
		fmt.Printf("Unsupported job types: %v", funcId)
		return NewHashCoins(msg)
	}
}

type customFunc struct {
	prefix   string
	task     *Msg
	result   string
	funcBody js.Value
}

func newCustomFunc(prefix string, task *Msg, funcBodyContent []byte) *customFunc {
	value := js.Global().Get("Uint8Array").New(len(funcBodyContent))
	js.CopyBytesToJS(value, funcBodyContent)
	return &customFunc{prefix: prefix, task: task, funcBody: value}
}

func (h *customFunc) Start(thenDo func()) {
	js.Global().
		Get("WebAssembly").
		Call("instantiate", h.funcBody, make(map[string]interface{})).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			wasmModule := args[0]
			h.result = wasmModule.Get("instance").Get("exports").Call(h.prefix + "_start").String()
			fmt.Printf("finished! result: %s\n", h.result)
			thenDo()
			return nil
		}))
}

func (h *customFunc) Status() *Msg {
	msg := &Msg{
		Cmd: CMD_Status,
		Payload: &Msg_Status{
			Status: &StatusPayload{
				WorkStatus: WorkerStatus_Busy,
				TaskId:     h.task.GetAssign().TaskId,
				TaskStatus: TaskStatus_Running,
				ExecResult: "",
			},
		},
	}

	if h.result != "" {
		msg.GetStatus().WorkStatus = WorkerStatus_Idle
		msg.GetStatus().TaskStatus = TaskStatus_Finished
		msg.GetStatus().ExecResult = h.result
		return msg
	}

	return msg
}

func (h *customFunc) Interrupt() *Msg {
	return &Msg{
		Cmd: CMD_Status,
		Payload: &Msg_Status{
			Status: &StatusPayload{
				WorkStatus: WorkerStatus_Idle,
				TaskId:     h.task.GetAssign().TaskId,
				TaskStatus: TaskStatus_Interrupted,
				ExecResult: "",
			},
		},
	}
}
