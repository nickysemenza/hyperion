package light

import (
	"time"

	"github.com/heatxsink/go-hue/lights"
)

type Bridge struct {
	Hostname string
	Username string
}

func (br *Bridge) SetColor(lightID int, color RGBColor, timing time.Duration) {
	x, y, _ := color.GetXyy()
	state := &lights.State{
		XY:             []float32{float32(x), float32(y)},
		Bri:            255,
		On:             true,
		TransitionTime: getTransitionTimeAs100msMultiple(timing),
	}

	hueLights := lights.New(br.Hostname, br.Username)
	hueLights.SetLightState(lightID, *state)
}

func getTransitionTimeAs100msMultiple(t time.Duration) uint16 {
	timingMs := uint16(t / time.Millisecond)
	return timingMs / 100
}
