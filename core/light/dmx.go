package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

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
}

type dmxOperation struct {
	universe, channel, value int
}

func (d *DMXLight) getChannelIDForAttributes(ctx context.Context, attrs ...string) (ids []int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getChannelIDForAttribute")
	defer span.Finish()
	profileMap := mainConfig.GetServerConfig(ctx).DMXProfiles
	profile, ok := profileMap[d.Profile]
	ids = make([]int, len(attrs))
	if ok {
		for x, attr := range attrs {
			channelIndex := getChannelIndexForAttribute(&profile, attr) //1 indexed
			ids[x] = d.StartAddress + channelIndex - 1
		}
		return
	}
	log.WithFields(log.Fields{"light": d.Name}).Warn("could not find DMX profile")
	return
}

//GetName returns the light's name.
func (d *DMXLight) GetName() string {
	return d.Name
}

//GetType returns the type of light.
func (d *DMXLight) GetType() string {
	return TypeDMX
}

//for a given color, blindly set the r,g, and b channels to that color, and update the state to reflect
func (d *DMXLight) blindlySetRGBToStateAndDMX(ctx context.Context, color color.RGB) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DMXLight blindlySetRGBToStateAndDMX")
	span.SetTag("dmx-name", d.Name)
	defer span.Finish()
	span.LogKV("event", "getting channel ids")
	rgbChannelIds := d.getChannelIDForAttributes(ctx, ChannelRed, ChannelGreen, ChannelBlue)
	span.LogKV("event", "getting channel values")
	rVal, gVal, bVal := color.AsComponents()

	span.LogKV("event", "begin getDMXStateInstance")
	ds := InitializeDMXState()
	span.LogKV("event", "now setting values")
	ds.set(ctx, dmxOperation{universe: d.Universe, channel: rgbChannelIds[0], value: rVal},
		dmxOperation{universe: d.Universe, channel: rgbChannelIds[1], value: gVal},
		dmxOperation{universe: d.Universe, channel: rgbChannelIds[2], value: bVal})

	SetCurrentState(d.Name, State{RGB: color})

}

//SetState updates the light's state.
//TODO: other properties? on/off?
func (d *DMXLight) SetState(ctx context.Context, target TargetState) {
	tickIntervalFadeInterpolation := mainConfig.GetServerConfig(ctx).Timings.FadeInterpolationTick
	currentState := GetCurrentState(d.Name)
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
	SetCurrentState(d.Name, target.ToState())

}

//DMXState holds the DMX512 values for each channel
type DMXState struct {
	universes map[int][]byte
	m         sync.Mutex
}

var (
	dmxStateInstance *DMXState
	once             sync.Once
)

//InitializeDMXState makes a singleton for dmxState
func InitializeDMXState() *DMXState {
	once.Do(func() {
		m := make(map[int][]byte)
		dmxStateInstance = &DMXState{universes: m}

	})
	return dmxStateInstance
}
func (s *DMXState) getValue(universe, channel int) int {
	return int(s.universes[universe][channel-1])
}

func (s *DMXState) set(ctx context.Context, ops ...dmxOperation) error {
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

func (s *DMXState) initializeUniverse(universe int) {
	if s.universes[universe] == nil {
		chans := make([]byte, 255)
		s.universes[universe] = chans
	}
}

//OLAClient is the interface for communicating with ola
type OLAClient interface {
	SendDmx(universe int, values []byte) (status bool, err error)
	Close()
}

//SendDMXWorker sends OLA the current dmxState across all universes
func SendDMXWorker(ctx context.Context, client OLAClient, tick time.Duration, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer client.Close()

	t := time.NewTimer(tick)
	defer t.Stop()
	log.Printf("timer started at %v", time.Now())

	for {
		select {
		case <-ctx.Done():
			log.Println("SendDMXWorker shutdown")
			return ctx.Err()
		case <-t.C:
			for k, v := range InitializeDMXState().universes {
				timer := prometheus.NewTimer(metrics.ExternalResponseTime.WithLabelValues("ola"))
				client.SendDmx(k, v)
				timer.ObserveDuration()
			}
			t.Reset(tick)
		}
	}
}

func getChannelIndexForAttribute(p *mainConfig.LightProfileDMX, attrName string) int {
	id, ok := p.Channels[attrName]
	if ok {
		return id
	}
	return 0
}
