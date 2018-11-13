package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/heatxsink/go-hue/hue"
	"github.com/heatxsink/go-hue/lights"
	mainConfig "github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
)

//Light is a light
type Light interface {
	GetName() string
	GetType() string
	SetState(context.Context, *Manager, TargetState)
}

//constants for the different types of lights
const (
	TypeHue     = "hue"
	TypeDMX     = "DMX"
	TypeGeneric = "generic"
)

type ctxKey int

//Context keys
const (
	ContextKeyLightNames ctxKey = iota
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

//Manager holds the state of lights
type Manager struct {
	byName        StateMap
	stateMapLock  sync.RWMutex
	hueConnection HueConnection
}

//SetCurrentState will set the current state for a light
func (m *Manager) SetCurrentState(name string, new State) {
	m.stateMapLock.Lock()
	defer m.stateMapLock.Unlock()
	m.byName[name] = new
}

//GetLightNames returns all the light names
func (m *Manager) GetLightNames() []string {
	m.stateMapLock.RLock()
	defer m.stateMapLock.RUnlock()
	keys := make([]string, 0, len(m.byName))
	for k := range m.byName {
		keys = append(keys, k)
	}
	return keys
}

//GetCurrentState will get the current state for a light
func (m *Manager) GetCurrentState(name string) *State {
	m.stateMapLock.RLock()
	defer m.stateMapLock.RUnlock()
	state, ok := m.byName[name]
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

//HueConnection represents a connection to a hue bridge
type HueConnection interface {
	SetLightState(lightID int, state lights.State) ([]hue.ApiResponse, error)
	GetAllLights() ([]lights.Light, error)
}

//Initialize parses light config
func Initialize(ctx context.Context, h HueConnection) (*Manager, error) {
	config := mainConfig.GetServerConfig(ctx)
	s := Manager{
		hueConnection: h,
	}
	s.byName = make(StateMap)

	ByName = make(NameMap)
	for i := range config.Lights.Hue {
		x := &config.Lights.Hue[i]
		ByName[x.Name] = &HueLight{
			HueID: x.HueID,
			Name:  x.Name,
		}
		s.SetCurrentState(x.Name, State{})
	}
	for i := range config.Lights.DMX {
		x := &config.Lights.DMX[i]
		ByName[x.Name] = &DMXLight{
			Name:         x.Name,
			StartAddress: x.StartAddress,
			Universe:     x.Universe,
			Profile:      x.Profile,
		}
		s.SetCurrentState(x.Name, State{})
	}
	for i := range config.Lights.Generic {
		x := &config.Lights.Generic[i]
		ByName[x.Name] = &GenericLight{
			Name: x.Name,
		}
		s.SetCurrentState(x.Name, State{})
	}

	return &s, nil
}
