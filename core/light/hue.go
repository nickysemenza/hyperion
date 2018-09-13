package light

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/heatxsink/go-hue/lights"
	mainConfig "github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/nickysemenza/hyperion/util/metrics"
)

//HueLight is a philips hue light.
type HueLight struct {
	HueID int    `json:"hue_id" yaml:"hue_id"`
	Name  string `json:"name" yaml:"name"`
	State State  `json:"state" yaml:"state"`
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

//GetState returns the light's state.
func (hl *HueLight) GetState() *State {
	return &hl.State
}

//SetState updates the Hue's state.
func (hl *HueLight) SetState(ctx context.Context, s State) {
	metrics.SetGagueWithNsFromTime(time.Now(), metrics.ResponseTimeNsHue)
	hl.m.Lock()
	defer hl.m.Unlock()

	hl.State = s
	go hl.SetColor(ctx, s.RGB, s.Duration)
}

//SetColor calls the Hue HTTP API to set the light's state to the given color, with given transition time (full brightness)
func (hl *HueLight) SetColor(ctx context.Context, color color.RGB, timing time.Duration) {
	lightID := hl.HueID
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

	log.WithFields(log.Fields{"hue_light_id": lightID, "timing": timing, "now": time.Now(), "brightness": brightness, "on": isOn}).Infof("HueLight SetColor: %s", color.TermString())
	hueConfig := mainConfig.GetServerConfig(ctx).Outputs.Hue
	hueLights := lights.New(hueConfig.Address, hueConfig.Username)
	hueLights.SetLightState(lightID, *state)
}

func getTransitionTimeAs100msMultiple(t time.Duration) uint16 {
	timingMs := uint16(t / time.Millisecond)
	return timingMs / 100
}
