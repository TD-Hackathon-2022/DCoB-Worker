package controller

import "syscall/js"

type eventCallBack func(this js.Value, args []js.Value) interface{}

func onOpen() eventCallBack {

	return func(this js.Value, args []js.Value) interface{} {
		ws
	}
}
