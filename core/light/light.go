package light

import (
	"context"
	"fmt"
	"time"

	"github.com/nickysemenza/hyperion/util/color"
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

//DebugString gives info
func DebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.GetName(), l.GetType())
}

//NameMap holds string-keyed Lights
type NameMap map[string]Light

//GetLights returns lights keyed by name
func GetLights() NameMap {
	return ByName
}

//ByName holds a name-keyed map of Lights
var ByName NameMap

//GetByName looks up a light by name
func GetByName(name string) Light {
	for _, x := range ByName {
		if x.GetName() == name {
			return x
		}
	}
	return nil
}
