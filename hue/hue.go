package hue

import (
	"time"

	"github.com/heatxsink/go-hue/lights"
	"github.com/lucasb-eyer/go-colorful"
)

type Bridge struct {
	Hostname string
	Username string
}

func (br *Bridge) SetRGB(lightID, r, g, b int) {
	cc := colorful.Color{
		R: float64(r / 255),
		G: float64(g / 255),
		B: float64(b / 255),
	}
	// x, y, _ := cc.Xyz()
	x, y, _ := colorful.XyzToXyy(colorful.LinearRgbToXyz(cc.LinearRgb()))
	state := &lights.State{
		XY:             []float32{float32(x), float32(y)},
		Bri:            255,
		On:             true,
		TransitionTime: 0,
	}

	hueLights := lights.New(br.Hostname, br.Username)
	hueLights.SetLightState(lightID, *state)
}

func Testing() {
	br := Bridge{
		Hostname: "10.0.1.55",
		Username: "alW0LsA1mnXB28T4txGs01BeHi1WBr661VZ1eqEF",
	}

	br.SetRGB(2, 255, 0, 0)
	time.Sleep(2 * time.Second)
	br.SetRGB(2, 0, 255, 0)
	time.Sleep(2 * time.Second)
	br.SetRGB(2, 0, 0, 255)
}
