package executor

import (
	"fmt"
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
)

type Executor interface {
	Start()
	Status() *Msg
	Interrupt() *Msg
}

func ExecutorBuilder(msg *Msg) Executor {
	if msg.GetAssign().FuncId != "hash-miner" {
		fmt.Printf("Unsupported job types: %v", msg.GetAssign().FuncId)
	}
	return NewHashCoins(msg)
}
