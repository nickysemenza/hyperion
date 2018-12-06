package light

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	mainConfig "github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
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
func (d *DMXLight) blindlySetRGBToStateAndDMX(ctx context.Context, m *Manager, color color.RGB) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DMXLight blindlySetRGBToStateAndDMX")
	span.SetTag("dmx-name", d.Name)
	defer span.Finish()
	span.LogKV("event", "getting channel ids")
	rgbChannelIds := d.getChannelIDForAttributes(ctx, ChannelRed, ChannelGreen, ChannelBlue)
	span.LogKV("event", "getting channel values")
	rVal, gVal, bVal := color.AsComponents()

	span.LogKV("event", "begin getDMXStateInstance")
	span.LogKV("event", "now setting values")
	m.dmxState.set(ctx, dmxOperation{universe: d.Universe, channel: rgbChannelIds[0], value: rVal},
		dmxOperation{universe: d.Universe, channel: rgbChannelIds[1], value: gVal},
		dmxOperation{universe: d.Universe, channel: rgbChannelIds[2], value: bVal})

	m.SetState(d.Name, State{RGB: color})

}

//SetState updates the light's state.
//TODO: other properties? on/off?
func (d *DMXLight) SetState(ctx context.Context, m *Manager, target TargetState) {
	tickIntervalFadeInterpolation := mainConfig.GetServerConfig(ctx).Timings.FadeInterpolationTick
	currentState := m.GetState(d.Name)
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
		d.blindlySetRGBToStateAndDMX(ctx, m, interpolated)

		time.Sleep(tickIntervalFadeInterpolation)
	}

	d.blindlySetRGBToStateAndDMX(ctx, m, target.RGB)
	span.LogKV("event", "finished fade interpolation")
	m.SetState(d.Name, target.ToState())

}

func getChannelIndexForAttribute(p *mainConfig.LightProfileDMX, attrName string) int {
	id, ok := p.Channels[attrName]
	if ok {
		return id
	}
	return 0
}
