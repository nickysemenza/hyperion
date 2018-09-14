package light

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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

	tickIntervalFadeInterpolation = time.Millisecond * 30
	tickIntervalSendToOLA         = time.Millisecond * 50
)

//DMXLight is a DMX light
type DMXLight struct {
	Name         string `json:"name" yaml:"name"`
	StartAddress int    `json:"start_address" yaml:"start_address"`
	Universe     int    `json:"universe" yaml:"universe"`
	Profile      string `json:"profile" yaml:"profile"`
	State        State  `json:"state" yaml:"state"`
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
	rChan, gChan, bChan := d.getRGBChannelIDs()
	rVal, gVal, bVal := color.AsComponents()

	ds := getDMXStateInstance()
	ds.setDMXValue(ctx, d.Universe, rChan, rVal)
	ds.setDMXValue(ctx, d.Universe, gChan, gVal)
	ds.setDMXValue(ctx, d.Universe, bChan, bVal)

	d.State.RGB = color

}

//SetState updates the light's state.
//TODO: other properties? on/off?
func (d *DMXLight) SetState(ctx context.Context, target State) {
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
		d.blindlySetRGBToStateAndDMX(ctx, color.GetRGBFromColorful(interpolated))

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

func (s *dmxState) setDMXValue(ctx context.Context, universe, channel, value int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "setDMXValue")
	defer span.Finish()
	span.SetTag("universe", universe)
	span.SetTag("channel", channel)
	span.SetTag("value", value)
	if channel < 1 || channel > 255 {
		return fmt.Errorf("dmx channel (%d) not in range", channel)
	}
	s.m.Lock()
	defer s.m.Unlock()
	s.initializeUniverse(universe)
	s.universes[universe][channel-1] = byte(value)
	return nil
}

func (s *dmxState) initializeUniverse(universe int) {
	u := s.universes[universe]
	if u == nil {
		chans := make([]byte, 255)
		s.universes[universe] = chans
	}
}

//SendDMXWorker sends OLA the current dmxState across all universes
func SendDMXWorker(ctx context.Context) {
	olaAddress := mainConfig.GetServerConfig(ctx).Outputs.OLA.Address
	//TODO: put this on a timer
	client := gola.New(olaAddress)
	defer client.Close()

	s := getDMXStateInstance()

	metrics.SetGagueWithNsFromTime(time.Now(), metrics.ResponseTimeNsOLA)
	for {
		for k, v := range s.universes {
			client.SendDmx(k, v)
		}
		time.Sleep(tickIntervalSendToOLA)
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
