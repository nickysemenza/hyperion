package light

import (
	"time"

	"github.com/heatxsink/go-hue/lights"
)

//HueLight is a philips hue light.
type HueLight struct {
	HueID int    `json:"hue_id"`
	Name  string `json:"name"`
	State State  `json:"state"`
	Test  int
}

//GetName returns the light's name.
func (hl *HueLight) GetName() string {
	return hl.Name
}

//GetType returns the type of light.
func (hl *HueLight) GetType() string {
	return TypeHue
}

//SetState updates the Hue's state.
func (hl *HueLight) SetState(s State) {
	hl.State = s
	hl.Test++
	Config.HueBridge.SetColor(hl.HueID, s.RGB, s.Duration)
}

//HueBridge holds credentials for communicating with hues.
type HueBridge struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
}

//SetColor calls the Hue HTTP API to set the light's state to the given color, with given transition time (full brightness)
func (br *HueBridge) SetColor(lightID int, color RGBColor, timing time.Duration) {
	x, y, _ := color.GetXyy()
	state := &lights.State{
		XY:             []float32{float32(x), float32(y)},
		Bri:            255, //TODO: set this...
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
