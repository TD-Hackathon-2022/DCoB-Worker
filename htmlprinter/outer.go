package htmlprinter

import "syscall/js"

func PrintHHtml(out string) {
	js.Global().Get("document").Call("getElementById", "tasklist").Call("append", out+"\n")
}

func AppendHtml(out string) {
	js.Global().Get("document").Call("getElementById", "taskstatus").Call("append", out+"\n")
}
