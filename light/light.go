package light

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Light is a light
type Light interface {
	SetColor(RGBColor)
	// getColor() string
	getType() string
	getName() string
}

//GetLightDebugString gives info
func GetLightDebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.getName(), l.getType())
}

type LightConfig struct {
	Loaded bool `json:"is_loaded"`
	Lights struct {
		Hue     []HueLight     `json:"hue"`
		Dmx     []DMXLight     `json:"dmx"`
		Generic []GenericLight `json:"generic"`
	} `json:"lights"`
	HueBridge Bridge `json:"hue"`
	Ola       struct {
		Hostname string `json:"hostname"`
	} `json:"ola"`
}

var Config LightConfig

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

	Config = *lc
	return *lc
}

func GetLightByName(name string) Light {
	//TODO: interate over config.Lights using reflect
	for _, x := range Config.Lights.Hue {
		if x.getName() == name {
			return &x
		}
	}
	for _, x := range Config.Lights.Dmx {
		if x.getName() == name {
			return &x
		}
	}
	for _, x := range Config.Lights.Generic {
		if x.getName() == name {
			return &x
		}
	}
	return nil
}
