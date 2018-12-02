package light

import (
	"context"
	"testing"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestGetType(t *testing.T) {
	tests := []struct {
		input    Light
		expected string
	}{
		{&DMXLight{}, TypeDMX},
		{&HueLight{}, TypeHue},
		{&GenericLight{}, TypeGeneric},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.input.GetType())
	}
}
func TestGenericLightInterface(t *testing.T) {
	l := &GenericLight{Name: "a"}
	s := DebugString(l)
	expected := "a - generic"
	if s != expected {
		t.Errorf("got %s, expected %s", s, expected)
	}

}

func TestDuplicateLights(t *testing.T) {
	s := &config.Server{
		Lights: config.Lights{
			Hue: []config.LightHue{{Name: "light1", HueID: 1}},
			DMX: []config.LightDMX{{Name: "light1"}},
		},
	}
	_, err := NewManager(s.InjectIntoContext(context.Background()), nil)
	require.Error(t, err)
}
func TestFindLightByName(t *testing.T) {
	dmx1 := &DMXLight{Name: "dmx1"}
	hue1 := &HueLight{Name: "hue1"}

	s := &config.Server{
		Lights: config.Lights{
			Hue: []config.LightHue{{Name: "hue1", HueID: 1}},
			DMX: []config.LightDMX{{Name: "dmx1"}},
		},
	}
	m, err := NewManager(s.InjectIntoContext(context.Background()), nil)
	require.NoError(t, err)

	tt := []struct {
		nameToFind string
		expected   Light
	}{
		{"hue1", hue1},
		{"dmx1", dmx1},
		{"aaa", nil},
	}
	for _, x := range tt {
		res := m.GetByName(x.nameToFind)
		if res == nil && x.expected == nil {
			continue
		}
		if res.GetName() != x.expected.GetName() {
			t.Errorf("got %s, expected %s", res, x.expected)
		}
	}
}

func TestManagerState(t *testing.T) {
	require := require.New(t)
	s := &config.Server{}
	m, err := NewManager(s.InjectIntoContext(context.Background()), nil)
	require.NoError(err)

	newState := State{RGB: color.RGB{R: 255}}
	m.SetState("light1", newState)
	m.SetState("light2", newState)
	require.Nil(m.GetState("foo"))
	require.Equal(255, m.GetState("light1").RGB.R)
	require.Equal(2, len(*m.GetAllStates()))

}
