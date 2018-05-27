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

func (d *DMXLight) getType() string {
	return "DMX"
}

func (d *DMXLight) getName() string {
	return d.Name
}
func (d *DMXLight) SetColor(c RGBColor) {

}

func TestSendDmx() {
	client := gola.New(Config.Ola.Hostname)
	defer client.Close()

	if x, err := client.GetDmx(3); err != nil {
		log.Printf("GetDmx: 1: %v", err)
	} else {
		log.Printf("GetDmx: 1: %v", x.Data)
	}
}
