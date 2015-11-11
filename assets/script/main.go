// +build js

package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"time"
)

var jq = jquery.NewJQuery

func main() {
	js.Global.Get("setInterval").Invoke(incrementClock, 1000)
	incrementClock()
	jq("#clock").Call("removeClass", "loading")
}

func incrementClock() {
	ct := fmt.Sprintf("%9x", time.Now().Unix())

	jq("#clock .digit").Call("each", func(i int, x *js.Object) {
		// i'th digit
		dg := ct[i : i+1]

		// Digit currently displayed
		cur_dig := jq("div.current", x)
		cd := cur_dig.Text()

		if cd != dg {
			new_dig := jq(".next", x).SetText(dg)

			new_dig.Call("toggleClass", "current", true).Call("toggleClass", "next", false)
			cur_dig.Call("toggleClass", "prev", true).Call("toggleClass", "current", false)

			js.Global.Get("setTimeout").Invoke(func() {
				jq(".prev", x).Call("remove")
				jq(x).Call("append", jq("<div class=\"next\"></div>"))
			}, 500)
		}
		if i == 0 && dg != " " {
			// Happy new Epoch!
			jq(x).Call("toggleClass", "empty", false)
		}
	})
}
