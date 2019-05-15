package light

import (
	"context"
	"sync"
	"time"

	"go.opencensus.io/trace"

	"github.com/heatxsink/go-hue/lights"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/nickysemenza/hyperion/util/metrics"
	"github.com/nickysemenza/hyperion/util/tracing"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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

//GetID returns the type of light.
func (hl *HueLight) GetID() string {
	return string(hl.HueID)
}

//SetState updates the Hue's state.
func (hl *HueLight) SetState(ctx context.Context, m Manager, s TargetState) {
	ctx, span := trace.StartSpan(ctx, "HueLight SetState")
	defer span.End()
	setSpanMeta(span, hl)
	span.Annotate([]trace.Attribute{}, "acquiring lock")
	hl.m.Lock()
	defer hl.m.Unlock()
	span.Annotate([]trace.Attribute{}, "acquired lock")
	m.SetState(hl.Name, s.ToState())
	hl.setColor(ctx, m.GetHueConnection(), s.RGB, s.Duration)
}

//setColor calls the Hue HTTP API to set the light's state to the given color, with given transition time (full brightness)
func (hl *HueLight) setColor(ctx context.Context, bridgeConn HueConnection, color color.RGB, timing time.Duration) error {
	ctx, span := trace.StartSpan(ctx, "HueLight setColor")
	defer span.End()
	// setSpanMeta(span, hl)

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

	log.WithFields(log.Fields{"hue_light_id": lightID, "timing": timing, "now": time.Now(), "brightness": brightness, "on": isOn}).Infof("HueLight setColor: %s", color.TermString())
	span.Annotate([]trace.Attribute{}, "sending hue light change to bridge")
	timer := prometheus.NewTimer(metrics.ExternalResponseTime.WithLabelValues("hue"))
	_, err := bridgeConn.SetLightState(lightID, *state) //TODO: use response
	if err != nil {
		tracing.SetError(span, err)
		log.WithFields(log.Fields{"hue_light_id": lightID}).Infof("HueLight setColor failure: %v", err)
		return err
	}
	timer.ObserveDuration()
	span.Annotate([]trace.Attribute{}, "done")
	return nil
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
func (m *StateManager) GetDiscoveredHues() DiscoveredHues {
	lights, _ := m.hueConnection.GetAllLights()

	byName := make(map[string]int)
	for _, x := range lights {
		byName[x.Name] = x.ID
	}

	return DiscoveredHues{lights, byName}
}

//GetHueConnection returns the hue connection object
func (m *StateManager) GetHueConnection() HueConnection {
	return m.hueConnection
}
