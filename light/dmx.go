package light

import (
	"log"

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
}

func testSendDmx() {
	client := gola.New(Config.Ola.Hostname)
	defer client.Close()

	if x, err := client.GetDmx(3); err != nil {
		log.Printf("GetDmx: 1: %v", err)
	} else {
		log.Printf("GetDmx: 1: %v", x.Data)
	}
}
