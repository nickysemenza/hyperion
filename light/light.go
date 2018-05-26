package light

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Light is a light
type Light interface {
	// setColor(string)
	// getColor() string
	getType() string
	getName() string
}

//GetLightDebugString gives info
func GetLightDebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.getName(), l.getType())
}

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

type HueLight struct {
	HueID int    `json:"hue_id"`
	Name  string `json:"name"`
}

func (hl *HueLight) getName() string {
	return hl.Name
}

func (hl *HueLight) getType() string {
	return "hue"
}

//GenericLight is for testing
type GenericLight struct {
	Name  string `json:"name"`
	Color string
}

func (gl *GenericLight) getType() string {
	return "GenericLight"
}
func (gl *GenericLight) getName() string {
	return gl.Name
}
func (gl *GenericLight) setColor(color string) {
	gl.Color = color
}
func (gl *GenericLight) getColor() string {
	return gl.Color
}

type LightConfig struct {
	Loaded bool `json:"is_loaded"`
	Lights struct {
		Hue     []HueLight     `json:"hue"`
		Dmx     []DMXLight     `json:"dmx"`
		Generic []GenericLight `json:"generic"`
	} `json:"lights"`
}

func ReadLightConfigFromFile(file string) LightConfig {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	lc := &LightConfig{}
	err = json.Unmarshal(raw, &lc)
	if err != nil {
		fmt.Println(err)
	}
	lc.Loaded = true

	return *lc
}

func GetLightByName(config *LightConfig, name string) Light {
	//TODO: interate over config.Lights using reflect
	for _, x := range config.Lights.Hue {
		if x.getName() == name {
			return &x
		}
	}
	for _, x := range config.Lights.Dmx {
		if x.getName() == name {
			return &x
		}
	}
	for _, x := range config.Lights.Generic {
		if x.getName() == name {
			return &x
		}
	}
	return nil
}
