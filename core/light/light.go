package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	mainConfig "github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
)

//Light is a light
type Light interface {
	GetName() string
	GetType() string
	SetState(context.Context, TargetState)
}

//constants for the different types of lights
const (
	TypeHue     = "hue"
	TypeDMX     = "DMX"
	TypeGeneric = "generic"
)

//TargetState represents the state of a light, is source of truth
type TargetState struct {
	// On   bool
	State
	Duration time.Duration `json:"duration"` //time to transition to the new state
}

//ToState converts a TargetState to a State
func (t *TargetState) ToState() State {
	return State{RGB: t.RGB}
}

//State represents the current state of the light
type State struct {
	RGB color.RGB `json:"rgb"` //RGB color
}

//NameMap holds string-keyed Lights
type NameMap map[string]Light

//StateMap holds Global light state
type StateMap map[string]State

type stateManager struct {
	byName       StateMap
	stateMapLock sync.RWMutex
}

var states stateManager

//SetCurrentState will set the current state for a light
func SetCurrentState(name string, s State) {
	states.stateMapLock.Lock()
	defer states.stateMapLock.Unlock()
	if states.byName == nil {
		states.byName = make(StateMap)
	}
	states.byName[name] = s
}

//GetLightNames returns all the light names
//TODO: move this to pull from config in context
func GetLightNames() []string {
	states.stateMapLock.RLock()
	defer states.stateMapLock.RUnlock()
	keys := make([]string, 0, len(states.byName))
	for k := range states.byName {
		keys = append(keys, k)
	}
	return keys
}

//GetCurrentState will get the current state for a light
func GetCurrentState(name string) *State {
	states.stateMapLock.RLock()
	defer states.stateMapLock.RUnlock()
	state, ok := states.byName[name]
	if ok {
		return &state
	}
	return nil
}

func (t *TargetState) String() string {
	return fmt.Sprintf("Duration: %s, RGB: %s", t.Duration, t.RGB.TermString())
}

//DebugString gives info
func DebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.GetName(), l.GetType())
}

//GetLightsByName returns lights keyed by name
func GetLightsByName() NameMap {
	return ByName
}

//ByName holds a name-keyed map of Lights
var ByName NameMap

//GetByName looks up a light by name
func GetByName(name string) Light {
	light, ok := ByName[name]
	if ok {
		return light
	}
	return nil
}

//Initialize parses light config
func Initialize(ctx context.Context) error {
	config := mainConfig.GetServerConfig(ctx)

	ByName = make(NameMap)
	for i := range config.Lights.Hue {
		x := &config.Lights.Hue[i]
		ByName[x.Name] = &HueLight{
			HueID: x.HueID,
			Name:  x.Name,
		}
		SetCurrentState(x.Name, State{})
	}
	for i := range config.Lights.DMX {
		x := &config.Lights.DMX[i]
		ByName[x.Name] = &DMXLight{
			Name:         x.Name,
			StartAddress: x.StartAddress,
			Universe:     x.Universe,
			Profile:      x.Profile,
		}
		SetCurrentState(x.Name, State{})
	}
	for i := range config.Lights.Generic {
		x := &config.Lights.Generic[i]
		ByName[x.Name] = &GenericLight{
			Name: x.Name,
		}
		SetCurrentState(x.Name, State{})
	}

	return nil
}
