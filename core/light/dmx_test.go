package light

import (
	"context"
	"testing"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/require"
)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DMXLight{
				StartAddress: tt.fields.StartAddress,
				Universe:     tt.fields.Universe,
				Profile:      tt.fields.Profile,
			}
			sm, _ := NewManager(ctx, nil)
			sm.dmxState = DMXState{universes: make(map[int][]byte)}
			d.blindlySetRGBToStateAndDMX(ctx, sm, tt.color)
			//green means first chan should be 0, secnd 255
			require.Equal(t, 0, sm.dmxState.getValue(tt.fields.Universe, tt.fields.StartAddress))
			require.Equal(t, 255, sm.dmxState.getValue(tt.fields.Universe, 1+tt.fields.StartAddress))
		})
	}
}
