package light

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nickysemenza/hyperion/backend/color"
)

//Light is a light
type Light interface {
	GetName() string
	GetType() string
	SetState(context.Context, State)
}

//Wrapper holds Lights
type Wrapper struct {
	Light Light `json:"light"`
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

//WrapperMap is a map
type WrapperMap map[string]Wrapper

//DMXProfileMap is a map of profiles
type DMXProfileMap map[string]dmxProfile

//GetWrapperMap returns lights keyed by name
func GetWrapperMap() WrapperMap {
	return ByName
}

//config is a global var containing the current lights
var config Config

//ByName holds a name-keyed map of LightWrappers
var ByName WrapperMap

//DMXProfilesByName holds dmx profiles
var DMXProfilesByName DMXProfileMap

//ReadLightConfigFromFile reads a config.json
func ReadLightConfigFromFile(file string) Config {
	//read file
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
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
	ByName = make(map[string]Wrapper)
	for i := range config.Lights.Hue {
		h := &config.Lights.Hue[i]
		ByName[h.GetName()] = Wrapper{h}
	}
	for i, x := range config.Lights.Dmx {
		ByName[x.GetName()] = Wrapper{&config.Lights.Dmx[i]}
	}
	for i, x := range config.Lights.Generic {
		ByName[x.GetName()] = Wrapper{&config.Lights.Generic[i]}
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
		if x.Light.GetName() == name {
			return x.Light
		}
	}
	return nil
}
