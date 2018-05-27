package light

import (
	"testing"
)

func TestGenericLightInterface(t *testing.T) {
	l := &GenericLight{Name: "a"}
	s := GetLightDebugString(l)
	expected := "a - generic"
	if s != expected {
		t.Errorf("got %s, expected %s", s, expected)
	}

}

func TestFindLightByName(t *testing.T) {
	lc := &Inventory{}
	generic1 := &GenericLight{Name: "test1"}
	dmx1 := &DMXLight{Name: "dmx1"}
	hue1 := &HueLight{Name: "hue1"}
	lc.Lights.Generic = []GenericLight{*generic1}
	lc.Lights.Hue = []HueLight{*hue1}
	lc.Lights.Dmx = []DMXLight{*dmx1}

	Config = *lc
	tt := []struct {
		nameToFind string
		expected   Light
	}{
		{"test1", generic1},
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
