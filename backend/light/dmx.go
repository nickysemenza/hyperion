package light

import (
	"fmt"
	"sync"

	"github.com/nickysemenza/gola"
)

//DMXLight is a DMX light
type DMXLight struct {
	Name         string `json:"name"`
	StartAddress int    `json:"start_address"`
	NumChannels  int    `json:"num_channels"`
	Universe     int    `json:"universe"`
	Profile      string `json:"profile"`
}

//GetName returns the light's name.
func (d *DMXLight) GetName() string {
	return d.Name
}

//GetType returns the type of light.
func (d *DMXLight) GetType() string {
	return TypeDMX
}

//SetState updates the light's state.
func (d *DMXLight) SetState(s State) {
	//TODO:
	//	get r/g/b channels based on mapping
	//	call setDMXValue on given channels
	//	other properties? on/off?
}

//dmxState holds the DMX512 values for each channel
type dmxState struct {
	universes map[int32][]byte
}

var dmxStateInstance *dmxState
var once sync.Once

//getDMXStateInstance makes a singleton for dmxState
func getDMXStateInstance() *dmxState {
	once.Do(func() {
		m := make(map[int32][]byte)
		dmxStateInstance = &dmxState{universes: m}

	})
	return dmxStateInstance
}

func (s *dmxState) setDMXValue(universe, channel, value int32) error {
	if channel < 1 || channel > 255 {
		return fmt.Errorf("dmx channel (%d) not in range", channel)
	}
	s.initializeUniverse(universe)
	s.universes[universe][channel-1] = byte(value)
	return nil
}

func (s *dmxState) initializeUniverse(universe int32) {
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
	for k, v := range s.universes {
		client.SendDmx(int(k), v)
	}
}
