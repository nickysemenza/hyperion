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

func (t *TargetState) String() string {
	return fmt.Sprintf("Duration: %s, RGB: %s", t.Duration, t.RGB.TermString())
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
	states        StateMap
	items         NameMap
	stateLock     sync.RWMutex
	hueConnection HueConnection
}

//SetState will set the current state for a light
func (m *Manager) SetState(name string, new State) {
	m.stateLock.Lock()
	defer m.stateLock.Unlock()
	m.states[name] = new
}

//GetLightNames returns all the light names
func (m *Manager) GetLightNames() []string {
	keys := make([]string, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

//GetState will get the current state for a light
func (m *Manager) GetState(name string) *State {
	m.stateLock.RLock()
	defer m.stateLock.RUnlock()
	state, ok := m.states[name]
	if ok {
		return &state
	}
	return nil
}

//GetAllStates will get the current state for all lights
func (m *Manager) GetAllStates() *StateMap {
	return &m.states
}

//DebugString gives info
func DebugString(l Light) string {
	return fmt.Sprintf("%s - %s", l.GetName(), l.GetType())
}

//GetLightsByName returns lights keyed by name
func (m *Manager) GetLightsByName() NameMap {
	return m.items
}

//GetByName looks up a light by name
func (m *Manager) GetByName(name string) Light {
	light, ok := m.items[name]
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

//NewManager parses light config
func NewManager(ctx context.Context, h HueConnection) (*Manager, error) {
	config := mainConfig.GetServerConfig(ctx)
	m := Manager{
		hueConnection: h,
		states:        make(StateMap),
		items:         make(NameMap),
	}
	//populate with each type of light
	for i := range config.Lights.Hue {
		x := &config.Lights.Hue[i]
		if _, ok := m.items[x.Name]; ok {
			err := fmt.Errorf("duplicate lights found! name=%s", x.Name)
			return nil, err
		}
		m.items[x.Name] = &HueLight{
			HueID: x.HueID,
			Name:  x.Name,
		}
		m.SetState(x.Name, State{})
	}
	for i := range config.Lights.DMX {
		x := &config.Lights.DMX[i]
		if _, ok := m.items[x.Name]; ok {
			err := fmt.Errorf("duplicate lights found! name=%s", x.Name)
			return nil, err
		}
		m.items[x.Name] = &DMXLight{
			Name:         x.Name,
			StartAddress: x.StartAddress,
			Universe:     x.Universe,
			Profile:      x.Profile,
		}
		m.SetState(x.Name, State{})
	}
	for i := range config.Lights.Generic {
		x := &config.Lights.Generic[i]
		if _, ok := m.items[x.Name]; ok {
			err := fmt.Errorf("duplicate lights found! name=%s", x.Name)
			return nil, err
		}
		m.items[x.Name] = &GenericLight{
			Name: x.Name,
		}
		m.SetState(x.Name, State{})
	}

	return &m, nil
}
