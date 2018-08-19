package light

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nickysemenza/hyperion/util/color"
	"github.com/pkg/errors"
)

//Light is a light
type Light interface {
	GetName() string
	GetType() string
	SetState(context.Context, State)
}

//constants for the different types of lights
const (
	TypeHue     = "hue"
	TypeDMX     = "DMX"
	TypeGeneric = "generic"
)

//State represents the state of a light, is source of truth
type State struct {
	// On   bool
	RGB      color.RGB     `json:"rgb"`      //RGB color
	Duration time.Duration `json:"duration"` //time to transition to the new state
}

func (s *State) String() string {
	return fmt.Sprintf("Duration: %s, RGB: %s", s.Duration, s.RGB.String())
}

//GetLightDebugString gives info
func GetLightDebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.GetName(), l.GetType())
}

//Config holds config data
type Config struct {
	Loaded bool `json:"is_loaded"`
	Lights struct {
		Hue     []HueLight     `json:"hue"`
		Dmx     []DMXLight     `json:"dmx"`
		Generic []GenericLight `json:"generic"`
	} `json:"lights"`
	HueBridge HueBridge `json:"hue"`
	Ola       struct {
		Hostname string `json:"hostname"`
	} `json:"ola"`
	Profiles []dmxProfile `json:"profiles"`
}

//StringMap holds string-keyed Lights
type StringMap map[string]Light

//DMXProfileMap is a map of profiles
type DMXProfileMap map[string]dmxProfile

//GetLights returns lights keyed by name
func GetLights() StringMap {
	return ByName
}

//config is a global var containing the current lights
var config Config

//ByName holds a name-keyed map of Lights
var ByName StringMap

//DMXProfilesByName holds dmx profiles
var DMXProfilesByName DMXProfileMap

//ReadLightConfigFromFile reads a config.json
func ReadLightConfigFromFile(file string) Config {
	//read file
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not read config file"))
		os.Exit(1)
	}

	//unmarshall json
	err = json.Unmarshal(raw, &config)
	if err != nil {
		fmt.Println(err)
	}

	//parse dmx profiles
	DMXProfilesByName = make(map[string]dmxProfile)
	for _, item := range config.Profiles {
		DMXProfilesByName[item.Name] = item
	}

	//parse lights
	ByName = make(StringMap)
	for i := range config.Lights.Hue {
		x := &config.Lights.Hue[i]
		ByName[x.GetName()] = x
	}
	for i := range config.Lights.Dmx {
		x := &config.Lights.Dmx[i]
		ByName[x.GetName()] = x
	}
	for i := range config.Lights.Generic {
		x := &config.Lights.Generic[i]
		ByName[x.GetName()] = x
	}

	//done
	config.Loaded = true
	return config
}

//GetConfig gives the current Configuration
func GetConfig() *Config {
	if !config.Loaded {
		log.Fatal("light configuration isn't loaded!")
	}
	return &config
}

//GetByName looks up a light by name
func GetByName(name string) Light {
	for _, x := range ByName {
		if x.GetName() == name {
			return x
		}
	}
	return nil
}
