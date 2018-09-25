package light

import (
	"context"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
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
func (hl *HueLight) SetState(ctx context.Context, s TargetState) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HueLight SetState")
	defer span.Finish()
	span.SetTag("hue-id", hl.HueID)
	span.SetTag("hue-name", hl.GetName())
	span.LogKV("event", "acquiring lock")
	hl.m.Lock()
	defer hl.m.Unlock()
	span.LogKV("event", "acquired lock")
	SetCurrentState(hl.Name, s.ToState())
	go hl.SetColor(ctx, s.RGB, s.Duration) //todo: goroutine might be defeating purpose of lock??
}

//SetColor calls the Hue HTTP API to set the light's state to the given color, with given transition time (full brightness)
func (hl *HueLight) SetColor(ctx context.Context, color color.RGB, timing time.Duration) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HueLight SetColor")
	defer span.Finish()
	span.SetTag("hue-id", hl.HueID)
	span.SetTag("hue-name", hl.GetName())

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
	span.LogEventWithPayload("sending hue light change to bridge", state)
	timer := prometheus.NewTimer(metrics.ExternalResponseTime.WithLabelValues("hue"))
	hueLights.SetLightState(lightID, *state) //TODO: use response
	timer.ObserveDuration()
	span.LogKV("event", "done")
}

func getTransitionTimeAs100msMultiple(t time.Duration) uint16 {
	timingMs := uint16(t / time.Millisecond)
	return timingMs / 100
}

//DiscoveredHues is for wrapping hue info
type DiscoveredHues struct {
	AllLights []lights.Light `json:"all_lights"`
	ByName    map[string]int `json:"by_name"`
}

//GetDiscoveredHues finds all the hues on the network
func GetDiscoveredHues(ctx context.Context) DiscoveredHues {

	hueConfig := mainConfig.GetServerConfig(ctx).Outputs.Hue
	hueLights := lights.New(hueConfig.Address, hueConfig.Username)

	lights, _ := hueLights.GetAllLights()

	byName := make(map[string]int)
	for _, x := range lights {
		byName[x.Name] = x.ID
	}

	return DiscoveredHues{lights, byName}
}
