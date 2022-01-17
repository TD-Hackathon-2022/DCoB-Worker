package executor

import (
	. "github.com/TD-Hackathon-2022/DCoB-Scheduler/api"
)

type Executor interface {
	Start()
	Status() *Msg
	Interrupt() *Msg
}
