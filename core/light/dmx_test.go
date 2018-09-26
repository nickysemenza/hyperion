package light

import (
	"context"
	"testing"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/require"
)

func TestDMXAttributeChannels(t *testing.T) {
	tt := []struct {
		profile  config.LightProfileDMX
		name     string
		expected int
	}{
		{config.LightProfileDMX{Channels: map[string]int{"red": 1, "green": 2}}, "red", 1},
		{config.LightProfileDMX{Channels: map[string]int{"red": 1, "green": 2}}, "blue", 0},
	}
	for _, tc := range tt {
		res := getChannelIndexForAttribute(&tc.profile, tc.name)
		if res != tc.expected {
			t.Errorf("got channel index %d, expected %d", res, tc.expected)
		}
	}
}
func TestDMX(t *testing.T) {
	s1 := getDMXStateInstance()
	s1.setDMXValues(context.Background(), dmxOperation{2, 22, 40})

	s2 := getDMXStateInstance()
	require.EqualValues(t, 40, s2.universes[2][21], "didn't set DMX state instance properly")
	require.Error(t, s2.setDMXValues(context.Background(), dmxOperation{2, 0, 2}), "should not allow channel 0")
	require.Equal(t, s1, s2, "should be a singleton!")
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
		{"setLightToGreen", fields{Profile: "a", Universe: 4, StartAddress: 1}, color.GetRGBFromString("green")},
		{"withOffsetStartAddress", fields{Profile: "a", Universe: 4, StartAddress: 10}, color.GetRGBFromString("green")},
	}

	c := config.Server{}
	c.DMXProfiles = make(config.DMXProfileMap)
	c.DMXProfiles["a"] = config.LightProfileDMX{Name: "a", Channels: map[string]int{"red": 1, "green": 2, "blue": 3}}

	ctx := c.InjectIntoContext(context.Background())
	Initialize(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DMXLight{
				StartAddress: tt.fields.StartAddress,
				Universe:     tt.fields.Universe,
				Profile:      tt.fields.Profile,
			}
			d.blindlySetRGBToStateAndDMX(ctx, tt.color)
			ds := getDMXStateInstance()
			//green means first chan should be 0, secnd 255
			require.Equal(t, 0, ds.getDmxValue(tt.fields.Universe, tt.fields.StartAddress))
			require.Equal(t, 255, ds.getDmxValue(tt.fields.Universe, 1+tt.fields.StartAddress))
		})
	}
}
