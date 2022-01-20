package htmlprinter

import "syscall/js"

func PrintPHtml(out string) {
	outer := js.Global().Get("document").Call("createElement", "p")
	outer.Set("innerHTML", out)
	js.Global().Get("document").Call("getElementById", "test").Call("appendChild", outer)
}

func PrintHHtml(out string) {
	outer := js.Global().Get("document").Call("createElement", "h3")
	outer.Set("innerHTML", out)
	js.Global().Get("document").Call("getElementById", "test").Call("appendChild", outer)
}
