package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/nickysemenza/gola"
	mainConfig "github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/nickysemenza/hyperion/util/metrics"
	opentracing "github.com/opentracing/opentracing-go"
)

//Holds strings for the different channel types
const (
	ChannelRed   = "red"
	ChannelGreen = "green"
	ChannelBlue  = "blue"
)

//DMXLight is a DMX light
type DMXLight struct {
	Name         string `json:"name" yaml:"name"`
	StartAddress int    `json:"start_address" yaml:"start_address"`
	Universe     int    `json:"universe" yaml:"universe"`
	Profile      string `json:"profile" yaml:"profile"`
	State        State  `json:"state" yaml:"state"`
}

type dmxOperation struct {
	universe, channel, value int
}

func (d *DMXLight) getProfile() *dmxProfile {
	return getDMXProfileByName(d.Profile)
}

func (d *DMXLight) getChannelIDForAttribute(attr string) int {
	profile := d.getProfile()
	if profile == nil {
		log.Println("cannot find profile!")
	}
	channelIndex := profile.getChannelIndexForAttribute(attr)
	return d.StartAddress + channelIndex
}
func (d *DMXLight) getRGBChannelIDs() (int, int, int) {
	return d.getChannelIDForAttribute(ChannelRed),
		d.getChannelIDForAttribute(ChannelGreen),
		d.getChannelIDForAttribute(ChannelBlue)
}

//GetName returns the light's name.
func (d *DMXLight) GetName() string {
	return d.Name
}

//GetType returns the type of light.
func (d *DMXLight) GetType() string {
	return TypeDMX
}

//GetState returns the light's state.
func (d *DMXLight) GetState() *State {
	return &d.State
}

//for a given color, blindly set the r,g, and b channels to that color, and update the state to reflect
func (d *DMXLight) blindlySetRGBToStateAndDMX(ctx context.Context, color color.RGB) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DMXLight blindlySetRGBToStateAndDMX")
	span.SetTag("dmx-name", d.Name)
	defer span.Finish()
	span.LogKV("event", "getting channel ids")
	rChan, gChan, bChan := d.getRGBChannelIDs()
	span.LogKV("event", "getting channel valyes")
	rVal, gVal, bVal := color.AsComponents()

	span.LogKV("event", "begin getDMXStateInstance")
	ds := getDMXStateInstance()
	span.LogKV("event", "now setting values")
	ds.setDMXValues(ctx, dmxOperation{universe: d.Universe, channel: rChan, value: rVal},
		dmxOperation{universe: d.Universe, channel: gChan, value: gVal},
		dmxOperation{universe: d.Universe, channel: bChan, value: bVal})

	d.State.RGB = color

}

//SetState updates the light's state.
//TODO: other properties? on/off?
func (d *DMXLight) SetState(ctx context.Context, target State) {
	tickIntervalFadeInterpolation := mainConfig.GetServerConfig(ctx).Timings.FadeInterpolationTick
	currentState := d.State
	numSteps := int(target.Duration / tickIntervalFadeInterpolation)
	span, ctx := opentracing.StartSpanFromContext(ctx, "DMX SetState")
	defer span.Finish()
	span.SetTag("dmx-name", d.Name)
	span.SetTag("target-duration-ms", target.Duration)

	log.Printf("dmx fade [%s] to [%s] over %d steps", currentState.RGB.TermString(), target.String(), numSteps)

	span.LogKV("event", "begin fade interpolation")
	for x := 0; x < numSteps; x++ {
		interpolated := currentState.RGB.GetInterpolatedFade(target.RGB, x, numSteps)
		//keep state updated:
		d.blindlySetRGBToStateAndDMX(ctx, interpolated)

		time.Sleep(tickIntervalFadeInterpolation)
	}

	d.blindlySetRGBToStateAndDMX(ctx, target.RGB)
	span.LogKV("event", "finished fade interpolation")
	d.State = target

}

//dmxState holds the DMX512 values for each channel
type dmxState struct {
	universes map[int][]byte
	m         sync.Mutex
}

var dmxStateInstance *dmxState
var once sync.Once

//getDMXStateInstance makes a singleton for dmxState
func getDMXStateInstance() *dmxState {
	once.Do(func() {
		m := make(map[int][]byte)
		dmxStateInstance = &dmxState{universes: m}

	})
	return dmxStateInstance
}
func (s *dmxState) getDmxValue(universe, channel int) int {
	return int(s.universes[universe][channel-1])
}

func (s *dmxState) setDMXValues(ctx context.Context, ops ...dmxOperation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "setDMXValues")
	defer span.Finish()
	span.SetTag("operations", ops)
	s.m.Lock()
	defer s.m.Unlock()
	for _, op := range ops {
		channel := op.channel
		universe := op.universe
		value := op.value
		if channel < 1 || channel > 255 {
			return fmt.Errorf("dmx channel (%d) not in range, op=%v", channel, op)
		}

		s.initializeUniverse(universe)
		s.universes[universe][channel-1] = byte(value)
	}

	return nil
}

func (s *dmxState) initializeUniverse(universe int) {
	if s.universes[universe] == nil {
		chans := make([]byte, 255)
		s.universes[universe] = chans
	}
}

//SendDMXWorker sends OLA the current dmxState across all universes
func SendDMXWorker(ctx context.Context) {
	olaConfig := mainConfig.GetServerConfig(ctx).Outputs.OLA
	if !olaConfig.Enabled {
		log.Info("ola output is not enabled")
		return
	}
	//TODO: put this on a timer
	client := gola.New(olaConfig.Address)
	defer client.Close()

	s := getDMXStateInstance()

	for {
		for k, v := range s.universes {
			timer := prometheus.NewTimer(metrics.ExternalResponseTime.WithLabelValues("ola"))
			client.SendDmx(k, v)
			timer.ObserveDuration()
		}
		time.Sleep(olaConfig.Tick)
	}
}

type dmxProfile struct {
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
	Channels     []string `json:"channels"`
}

func (p *dmxProfile) getChannelIndexForAttribute(attrName string) int {

	for i, x := range p.Channels {
		if attrName == x {
			return i
		}
	}
	return -1
}

//DMXProfileMap is a map of profiles
type DMXProfileMap map[string]dmxProfile

//DMXProfilesByName holds dmx profiles
var DMXProfilesByName DMXProfileMap

func getDMXProfileByName(name string) *dmxProfile {
	for _, x := range DMXProfilesByName {
		if x.Name == name {
			return &x
		}
	}
	return nil
}
