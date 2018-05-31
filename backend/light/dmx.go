package light

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nickysemenza/gola"
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
	Name         string `json:"name"`
	StartAddress int    `json:"start_address"`
	NumChannels  int    `json:"num_channels"`
	Universe     int    `json:"universe"`
	Profile      string `json:"profile"`
	State        State  `json:"state"`
}

func (d *DMXLight) getChannelIDForAttribute(attr string) int {
	switch attr {
	case ChannelRed:
		return 1
	case ChannelGreen:
		return 2
	case ChannelBlue:
		return 3
	}
	return 0 //TODO: error
}
func (d *DMXLight) getRGBChannelIDs() (int, int, int) {
	//TODO: get channels based on profile
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

//for a given color, blindly set the r,g, and b channels to that color, and update the state to reflect
func (d *DMXLight) blindlySetRGBToStateAndDMX(color RGBColor) {
	rChan, gChan, bChan := d.getRGBChannelIDs()
	rVal, gVal, bVal := color.AsComponents()

	ds := getDMXStateInstance()
	ds.setDMXValue(d.Universe, rChan, rVal)
	ds.setDMXValue(d.Universe, gChan, gVal)
	ds.setDMXValue(d.Universe, bChan, bVal)

	d.State.RGB = color

}

//SetState updates the light's state.
//TODO: other properties? on/off?
func (d *DMXLight) SetState(target State) {
	currentState := d.State
	numSteps := int(target.Duration / tickIntervalFadeInterpolation)

	log.Printf("dmx fade [%s] to [%s] over %d steps", currentState.RGB.String(), target.String(), numSteps)

	c1 := currentState.RGB.asColorful()
	c2 := target.RGB.asColorful()

	for x := 1; x <= numSteps; x++ {
		interpolatedColor := c1.BlendHcl(c2, float64(x)/float64(numSteps)).Clamped()
		interpolatedRGB := getRGBFromColorful(interpolatedColor)

		//keep state updated:
		d.blindlySetRGBToStateAndDMX(interpolatedRGB)

		time.Sleep(tickIntervalFadeInterpolation)
	}

	d.blindlySetRGBToStateAndDMX(target.RGB)
	d.State = target

}

//dmxState holds the DMX512 values for each channel
type dmxState struct {
	universes map[int][]byte
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

func (s *dmxState) setDMXValue(universe, channel, value int) error {
	if channel < 1 || channel > 255 {
		return fmt.Errorf("dmx channel (%d) not in range", channel)
	}
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

//SendDMXValuesToOLA sends OLA the current dmxState across all universes
func SendDMXValuesToOLA() {
	//TODO: put this on a timer
	client := gola.New(Config.Ola.Hostname)
	defer client.Close()

	s := getDMXStateInstance()

	for {
		for k, v := range s.universes {
			client.SendDmx(k, v)
		}
		time.Sleep(tickIntervalSendToOLA)
	}
}
