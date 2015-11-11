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
		dg := ct[i : i+1]

		jq(x).SetText(dg).Call("toggleClass", "empty", dg == " ")
	})
}
