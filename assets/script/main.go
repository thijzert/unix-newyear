// +build js

package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"net/url"
	"time"
)

var jq = jquery.NewJQuery

func main() {
	js.Global.Get("setInterval").Invoke(incrementClock, 1000)
	incrementClock()
	jq("#clock").Call("removeClass", "loading")
}

func incrementClock() {
	now := time.Now().Unix()
	ct := fmt.Sprintf("%9x", now)

	hash := js.Global.Get("location").Get("hash").String()
	opts, _ := url.ParseQuery(hash[1:])
	fmt.Printf("%+v", opts)

	var animate = true
	if opts.Get("no-animation") != "" {
		animate = false
	}

	jq("#clock .digit").Call("each", func(i int, x *js.Object) {
		// i'th digit
		dg := ct[i : i+1]

		// Digit currently displayed
		cur_dig := jq("div.current", x)
		cd := cur_dig.Text()

		if cd != dg {
			if animate {
				new_dig := jq(".next", x).SetText(dg)

				new_dig.Call("toggleClass", "current", true).Call("toggleClass", "next", false)
				cur_dig.Call("toggleClass", "prev", true).Call("toggleClass", "current", false)

				js.Global.Get("setTimeout").Invoke(func() {
					jq(".prev", x).Call("remove")
					jq(x).Call("append", jq("<div class=\"next\"></div>"))
				}, 500)
			} else {
				cur_dig.SetText(dg)
			}
		}
		if i == 0 && dg != " " {
			// Happy new Epoch!
			jq(x).Call("toggleClass", "empty", false)
		}
	})
}
