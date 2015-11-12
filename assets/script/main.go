// +build js

package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"math"
	"math/rand"
	"net/url"
	"time"
)

const GRAVITY float32 = 0.05

var jq = jquery.NewJQuery
var tableau *js.Object
var fireworkShow []*Firework

func main() {
	tableau = js.Global.Get("document").Call("getElementById", "firework-canvas")
	js.Global.Get("setInterval").Invoke(incrementClock, 1000)
	incrementClock()
	jq("#clock").Call("removeClass", "loading")
}

func incrementClock() {
	now := int(time.Now().Unix())
	ct := fmt.Sprintf("%9x", now)

	hash := js.Global.Get("location").Get("hash").String()
	opts, _ := url.ParseQuery(hash[1:])

	var animate = true
	if opts.Get("no-animation") != "" {
		animate = false
	}

	var modulus int = 0x100000
	if opts.Get("wait-for") != "" {
		wait_for := opts.Get("wait-for")
		var nmod int = 0
		n, er := fmt.Sscanf(wait_for, "%x", &nmod)
		if n == 1 && nmod > 16 {
			modulus = nmod
		} else {
			fmt.Print(er)
		}
	}

	final_countdown := (modulus - (now % modulus)) % modulus

	if final_countdown == 0 {
		if fireworkShow == nil {
			fireworkShow = make([]*Firework, 3)
			for i, _ := range fireworkShow {
				fireworkShow[i] = NewFirework(1, 2, 25)
			}
		}
		for _, f := range fireworkShow {
			js.Global.Get("setTimeout").Invoke(func() {
				f.Start()
			}, 300+rand.Intn(1000))
		}
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

	if final_countdown < 16 {
		timeout := 700
		if final_countdown == 0 {
			timeout = 2100
		}
		jq("#final-countdown").SetText(fmt.Sprintf("%X", final_countdown)).Call("css", js.M{"color": "#f4f4f0", "opacity": 1.0}).Call("animate", js.M{"color": "#403040", "opacity": 0.0}, timeout)
	}
}

type Firework struct {
	X, Y, Step int
	Numsparks  int
	Sparks     []*Spark
	Animating  bool
}

func NewFirework(x, y, sparks int) *Firework {
	rv := &Firework{X: x, Y: y, Step: 0,
		Numsparks: sparks, Sparks: make([]*Spark, sparks)}

	for i, _ := range rv.Sparks {
		rv.Sparks[i] = NewSpark(3.0)
	}

	return rv
}

func (f *Firework) Explode() {
	x := float32(50 + rand.Intn(jq(tableau).Width()-200))
	y := float32(50 + rand.Intn(jq(tableau).Height()-200))
	colour := 1 + rand.Intn(7)
	style := rand.Intn(5)

	for _, d := range f.Sparks {
		d.Fire(x, y, colour, style)
	}
}

func (f *Firework) Animate() {
	if !f.Animating {
		return
	}

	js.Global.Get("setTimeout").Invoke(func() {
		f.Animate()
	}, 40)

	if f.Step > 70 {
		f.Step = 0
	}
	if f.Step == 0 {
		f.Explode()
	}

	f.Step++
	for _, d := range f.Sparks {
		d.Animate(f.Step)
	}
}

func (f *Firework) Start() {
	f.Animating = true
	for _, s := range f.Sparks {
		s.Box.Visible = true
	}
	f.Animate()
}

func (f *Firework) Stop() {
	f.Animating = false
	for _, s := range f.Sparks {
		s.Box.Visible = false
	}
}

type Spark struct {
	X, Y, Size    float32
	VX, VY        float32
	Colour, Style int
	Box           *DomBox
}

func NewSpark(size float32) *Spark {
	rv := &Spark{}
	rv.Box = NewDomBox(0, 0, size, size, "000", false)
	return rv
}

func (s *Spark) Fire(x, y float32, colour, style int) {
	s.X = x
	s.Y = y
	s.Colour = colour
	s.Style = style

	s.Box.X = x
	s.Box.Y = y
	s.Box.Visible = true
	s.Box.Restyle()

	var a, r float32
	a = rand.Float32() * 6.294
	r = 4.0

	if style == 0 {
		r = 2.0
		if rand.Float32() <= 0.6 {
			r = rand.Float32() * 2
		}
	} else if style == 1 {
		r = rand.Float32() * 2
	} else if style == 2 {
		r = 2.0
	} else if style == 3 {
		r = a - rand.Float32()
	} else if style == 4 {
		if rand.Float32() > 0.5 {
			s.Colour = 1
			r = 2.0 - 0.40*rand.Float32()
		} else {
			s.Colour = 7
			r = 1.0 - 0.25*rand.Float32()
		}
	}

	s.VX = r * float32(math.Sin(float64(a)))
	s.VY = r*float32(math.Cos(float64(a))) - 2
}

func (s *Spark) SetColour(intensity int) {
	if intensity < 0 {
		intensity = 0
	} else if intensity > 255 {
		intensity = 255
	}

	var css string
	if s.Colour == 1 {
		css = fmt.Sprintf("0000%2x", intensity)
	} else if s.Colour == 2 {
		css = fmt.Sprintf("00%2x00", intensity)
	} else if s.Colour == 4 {
		css = fmt.Sprintf("%2x0000", intensity)
	} else if s.Colour == 3 {
		css = fmt.Sprintf("00%2x%2x", intensity, intensity)
	} else if s.Colour == 5 {
		css = fmt.Sprintf("%2x00%2x", intensity, intensity)
	} else if s.Colour == 6 {
		css = fmt.Sprintf("%2x%2x00", intensity, intensity)
	} else if s.Colour == 7 {
		css = fmt.Sprintf("%2x%2x%2x", intensity, intensity, intensity)
	} else {
		ncl := 1 + rand.Intn(7)
		fmt.Printf("Unknown colour mode %d. Falling back to %d.", s.Colour, ncl)
		s.Colour = ncl
		s.SetColour(intensity)
		return
	}

	s.Box.Colour = css
}

func (s *Spark) Animate(step int) {
	var colour = 255 - (4 * step)
	if step > 30 {
		colour = int(rand.Float32() * float32(356-(step*4)))
	}
	s.SetColour(colour)

	s.VY += GRAVITY
	s.X += s.VX
	s.Y += s.VY
	s.Box.X = s.X
	s.Box.Y = s.Y
	s.Box.Restyle()
}

type DomBox struct {
	X, Y, W, H float32
	Colour     string
	Visible    bool
	domElt     *js.Object
}

func NewDomBox(x, y, w, h float32, colour string, visible bool) *DomBox {
	rv := &DomBox{X: x, Y: y, W: w, H: h, Colour: colour, Visible: visible}

	rv.domElt = js.Global.Get("document").Call("createElement", "div")
	tableau.Call("appendChild", rv.domElt)
	return rv
}

func (b *DomBox) Restyle() {
	disp := "block"
	if !b.Visible {
		disp = "none"
	}

	style := fmt.Sprintf("position:absolute;left:%gpx;top:%gpx;width:%gpx;height:%gpx;display:%s;background-color:#%s",
		b.X, b.Y, b.W, b.H, disp, b.Colour)

	b.domElt.Call("setAttribute", "style", style)
}
