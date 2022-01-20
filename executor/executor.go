package executor

import (
	"encoding/base64"
	"fmt"
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
	"strings"
	"syscall/js"
)

type Executor interface {
	Start()
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
		ch := loadCustomFunc(msg.GetAssign().Data, funcPrefixArr[1])
		return newCustomFunc(ch)
	}

	switch funcId {
	case hashMiner:
		return NewHashCoins(msg)
	default:
		fmt.Printf("Unsupported job types: %v", funcId)
		return NewHashCoins(msg)
	}
}

func loadCustomFunc(data string, prefix string) <-chan string {
	content, _ := base64.StdEncoding.DecodeString(data)
	prefixCh := make(chan string)
	value := js.Global().Get("Uint8Array").New(len(content))
	js.CopyBytesToJS(value, content)
	g := js.Global().Get("Go").New()
	js.Global().
		Get("WebAssembly").
		Call("instantiate", value, g.Get("importObject")).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			g.Call("run", args[0].Get("instance"))
			prefixCh <- prefix
			return nil
		}))

	return prefixCh
}

type customFunc struct {
	prefix string
}

func newCustomFunc(prefix <-chan string) *customFunc {
	return &customFunc{prefix: <-prefix}
}

func (h *customFunc) Start() {
	js.Global().Call(h.prefix + "_start")
}

func (h *customFunc) Status() *Msg {
	return nil
}

func (h *customFunc) Interrupt() *Msg {
	return nil
}
