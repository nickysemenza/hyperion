package light

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/heatxsink/go-hue/lights"
	"github.com/nickysemenza/hyperion/color"
)

//HueLight is a philips hue light.
type HueLight struct {
	HueID int    `json:"hue_id"`
	Name  string `json:"name"`
	State State  `json:"state"`
	m     sync.Mutex
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
func (hl *HueLight) SetState(ctx context.Context, s State) {
	hl.m.Lock()
	defer hl.m.Unlock()

	hl.State = s
	go GetConfig().HueBridge.SetColor(ctx, hl.HueID, s.RGB, s.Duration)
}

//HueBridge holds credentials for communicating with hues.
type HueBridge struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
}

//SetColor calls the Hue HTTP API to set the light's state to the given color, with given transition time (full brightness)
func (br *HueBridge) SetColor(ctx context.Context, lightID int, color color.RGB, timing time.Duration) {
	x, y, _ := color.GetXyy()
	brightness := uint8(255)
	isOn := true
	if color.IsBlack() {
		brightness = 0
		isOn = false
	}
	state := &lights.State{
		XY:             []float32{float32(x), float32(y)},
		Bri:            brightness,
		On:             isOn,
		TransitionTime: getTransitionTimeAs100msMultiple(timing),
	}

	log.WithFields(log.Fields{"hue_light_id": lightID, "timing": timing, "now": time.Now(), "brightness": brightness, "on": isOn}).Infof("HueLight SetColor: %s", color.String())
	hueLights := lights.New(br.Hostname, br.Username)
	hueLights.SetLightState(lightID, *state)
}

func getTransitionTimeAs100msMultiple(t time.Duration) uint16 {
	timingMs := uint16(t / time.Millisecond)
	return timingMs / 100
}
