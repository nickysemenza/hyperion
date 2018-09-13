package light

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

//Config holds config data
type Config struct {
	Loaded bool `json:"is_loaded" yaml:"is_loaded"`
	Lights struct {
		Hue     []HueLight     `json:"hue" yaml:"hue"`
		Dmx     []DMXLight     `json:"dmx" yaml:"dmx"`
		Generic []GenericLight `json:"generic" yaml:"generic"`
	} `json:"lights" yaml:"lights"`
	Profiles []dmxProfile `json:"profiles" yaml:"profiles"`
}

//config is a global var containing the current lights
var config Config

//ReadLightConfigFromFile reads a config.yaml
func ReadLightConfigFromFile(file string) Config {
	//read file
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not read config file"))
		os.Exit(1)
	}

	//unmarshall json
	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		fmt.Println(err)
	}

	//parse dmx profiles
	DMXProfilesByName = make(map[string]dmxProfile)
	for _, item := range config.Profiles {
		DMXProfilesByName[item.Name] = item
	}

	//parse lights
	ByName = make(NameMap)
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
