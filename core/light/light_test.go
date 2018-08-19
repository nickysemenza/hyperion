package light

import (
	"testing"

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

func TestFindLightByName(t *testing.T) {
	dmx1 := &DMXLight{Name: "dmx1"}
	hue1 := &HueLight{Name: "hue1"}

	ByName = make(NameMap)
	ByName["hue1"] = hue1
	ByName["dmx1"] = dmx1

	tt := []struct {
		nameToFind string
		expected   Light
	}{
		{"hue1", hue1},
		{"dmx1", dmx1},
		{"aaa", nil},
	}
	for _, x := range tt {
		res := GetByName(x.nameToFind)
		if res == nil && x.expected == nil {
			continue
		}
		if res.GetName() != x.expected.GetName() {
			t.Errorf("got %s, expected %s", res, x.expected)
		}
	}
}
