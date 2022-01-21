package executor

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"worker/htmlprinter"
)

type hashCoins struct {
	ctx    context.Context
	cancel context.CancelFunc
	sobel  func(pam ...interface{}) interface{}
}

func NewHashCoins(m *Msg) *hashCoins {
	ctx, cancel := context.WithCancel(context.Background())
	return &hashCoins{
		ctx:    context.WithValue(ctx, "task", m),
		cancel: cancel,
		sobel:  miningCoins,
	}
}

func (h *hashCoins) Start(thenDo func()) {
	defer h.cancel()
	assignTask, _ := h.ctx.Value("task").(*Msg)
	bit, _ := strconv.Atoi(assignTask.GetAssign().Data)
	data := h.sobel(bit).([]byte)
	sENC := base64.StdEncoding.EncodeToString(data)
	msg := &Msg{
		Cmd: CMD_Status,
		Payload: &Msg_Status{
			Status: &StatusPayload{
				WorkStatus: WorkerStatus_Idle,
				TaskId:     assignTask.GetAssign().TaskId,
				TaskStatus: TaskStatus_Finished,
				ExecResult: sENC,
			},
		},
	}
	h.ctx = context.WithValue(h.ctx, "result", msg)
	fmt.Printf("Result found: %v\n", msg)
	thenDo()
}

func (h *hashCoins) Status() *Msg {
	assignTask, _ := h.ctx.Value("task").(*Msg)
	msg := &Msg{
		Cmd: CMD_Status,
		Payload: &Msg_Status{
			Status: &StatusPayload{
				WorkStatus: WorkerStatus_Busy,
				TaskId:     assignTask.GetAssign().TaskId,
				TaskStatus: TaskStatus_Running,
				ExecResult: "",
			},
		},
	}

	if result, exist := h.ctx.Value("result").(*Msg); exist {
		return result
	}

	return msg
}

func (h *hashCoins) Interrupt() *Msg {
	assignTask, _ := h.ctx.Value("task").(*Msg)
	msg := &Msg{
		Cmd: CMD_Status,
		Payload: &Msg_Status{
			Status: &StatusPayload{
				WorkStatus: WorkerStatus_Idle,
				TaskId:     assignTask.GetAssign().TaskId,
				TaskStatus: TaskStatus_Interrupted,
				ExecResult: "",
			},
		},
	}
	defer h.cancel()
	return msg
}

func genSource() <-chan []byte {
	srcChan := make(chan []byte, 10)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for {
			b := make([]byte, 32)
			r.Read(b)
			srcChan <- sha256.New().Sum(b)
		}
	}()
	return srcChan
}

func miningCoins(pam ...interface{}) interface{} {
	srcChan := genSource()
	Ruler := strings.Repeat("0", pam[0].(int))
	coin := make([]byte, 32)
	for src := range srcChan {
		tmpStr := hex.EncodeToString(src)
		if strings.HasPrefix(tmpStr, Ruler) {
			time.Sleep(time.Second * 1)
			htmlprinter.AppendHtml("Coin found:" + tmpStr[:32])
			coin = src
			return coin
		}
	}
	panic("can't execute to this point")
}
