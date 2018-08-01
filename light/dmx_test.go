package light

import (
	"testing"

	"github.com/nickysemenza/hyperion/color"
	"github.com/stretchr/testify/assert"
)

func TestDMXAttributeChannels(t *testing.T) {
	tt := []struct {
		profile  dmxProfile
		name     string
		expected int
	}{
		{dmxProfile{Channels: []string{"red", "green"}}, "red", 0},
	}
	for _, tc := range tt {
		res := tc.profile.getChannelIndexForAttribute(tc.name)
		if res != tc.expected {
			t.Errorf("got channel index %d, expected %d", res, tc.expected)
		}
	}
}
func TestDMX(t *testing.T) {
	s1 := getDMXStateInstance()
	s1.setDMXValue(2, 22, 40)

	s2 := getDMXStateInstance()
	if s2.universes[2][22-1] != 40 {
		t.Error("didn't set DMX state instance properly")
	}

	if err := s2.setDMXValue(2, 0, 2); err == nil {
		t.Error("should not allow channel 0")
	}

	if s1 != s2 {
		t.Error("should be singleton!")
	}
}

func TestDMXLight_blindlySetRGBToStateAndDMX(t *testing.T) {
	type fields struct {
		StartAddress int
		Universe     int
		Profile      string
	}
	tests := []struct {
		name   string
		fields fields
		color  color.RGB
	}{
		{"setLightToGreen", fields{Profile: "a", Universe: 4}, color.GetRGBFromHex("#00FF00")},
		{"withOffsetStartAddress", fields{Profile: "a", Universe: 4, StartAddress: 10}, color.GetRGBFromHex("#00FF00")},
	}

	DMXProfilesByName = make(map[string]dmxProfile)
	DMXProfilesByName["a"] = dmxProfile{Name: "a", Channels: []string{"noop", "red", "green", "blue"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DMXLight{
				StartAddress: tt.fields.StartAddress,
				Universe:     tt.fields.Universe,
				Profile:      tt.fields.Profile,
			}
			d.blindlySetRGBToStateAndDMX(tt.color)
			ds := getDMXStateInstance()
			assert.Equal(t, 255, ds.getDmxValue(tt.fields.Universe, 2+tt.fields.StartAddress))
			assert.Equal(t, 0, ds.getDmxValue(tt.fields.Universe, 1+tt.fields.StartAddress))
		})
	}
}
