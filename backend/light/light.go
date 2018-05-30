package light

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Light is a light
type Light interface {
	GetName() string
	GetType() string
	SetState(State)
}

//Wrapper holds Lights
type Wrapper struct {
	Light Light
}

//constants for the different types of lights
const (
	TypeHue     = "hue"
	TypeDMX     = "DMX"
	TypeGeneric = "generic"
)

//GetLightDebugString gives info
func GetLightDebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.GetName(), l.GetType())
}

//Inventory holds config data
type Inventory struct {
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
}

//Config is a global var containing the current lights
var Config Inventory

//WrapperMap holds a name-keyed map of LightWrappers
var WrapperMap map[string]Wrapper

//ReadLightConfigFromFile reads a config.json
func ReadLightConfigFromFile(file string) Inventory {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	lc := &Inventory{}
	err = json.Unmarshal(raw, &lc)
	if err != nil {
		fmt.Println(err)
	}
	lc.Loaded = true

	WrapperMap = make(map[string]Wrapper)
	for i, x := range lc.Lights.Hue {
		WrapperMap[x.GetName()] = Wrapper{&lc.Lights.Hue[i]}
	}
	for i, x := range lc.Lights.Dmx {
		WrapperMap[x.GetName()] = Wrapper{&lc.Lights.Dmx[i]}
	}
	for i, x := range lc.Lights.Generic {
		WrapperMap[x.GetName()] = Wrapper{&lc.Lights.Generic[i]}
	}

	Config = *lc
	return *lc

}

//GetByName looks up a light by name
func GetByName(name string) Light {
	for _, x := range WrapperMap {
		if x.Light.GetName() == name {
			return x.Light
		}
	}
	return nil
}
