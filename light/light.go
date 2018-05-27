package light

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//RGBColor holds RGB values (0-255)
type RGBColor struct {
	R int
	G int
	B int
}

func (c *RGBColor) FancyString() string {

	red := color.New(color.BgRed).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	blue := color.New(color.BgBlue).SprintFunc()
	return fmt.Sprintf("%s %s %s", red(c.R), green(c.G), blue(c.B))
}

func (c *RGBColor) GetXyy() (x, y, Yout float64) {
	cc := colorful.Color{
		R: float64(c.R / 255),
		G: float64(c.G / 255),
		B: float64(c.B / 255),
	}
	//OLD:  x, y, _ := cc.Xyz()
	return colorful.XyzToXyy(colorful.LinearRgbToXyz(cc.LinearRgb()))
}

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
func (hl *HueLight) SetColor(c RGBColor) {
	Config.HueBridge.SetColor(hl.HueID, c, time.Duration(time.Second))
}

//GenericLight is for testing
type GenericLight struct {
	Name  string `json:"name"`
	Color RGBColor
}

func (gl *GenericLight) getType() string {
	return "GenericLight"
}
func (gl *GenericLight) getName() string {
	return gl.Name
}

func (gl *GenericLight) SetColor(c RGBColor) {
	gl.Color = c
}

type LightConfig struct {
	Loaded bool `json:"is_loaded"`
	Lights struct {
		Hue     []HueLight     `json:"hue"`
		Dmx     []DMXLight     `json:"dmx"`
		Generic []GenericLight `json:"generic"`
	} `json:"lights"`
	HueBridge Bridge `json:"hue"`
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
