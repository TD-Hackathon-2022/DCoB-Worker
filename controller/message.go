package controller

import "encoding/json"

type Message struct {
	CMD     int     `json:"cmd"`
	Payload PAYLOAD `json:"payload"`
}

type PAYLOAD struct {
	TaskId       string `json:"taskId,omitempty"`
	WorkerStatus string `json:"workerStatus,omitempty"`
	TaskStatus   int    `json:"taskStatus,omitempty"`
	ExecResult   string `json:"execResult,omitempty"`
	FuncId       string `json:"funcId,omitempty"`
}

func (p PAYLOAD) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (m *Message) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
