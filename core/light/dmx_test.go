package light

import (
	"context"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/mock"
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
			d.blindlySetRGBToStateAndDMX(ctx, sm, tt.color)
			//green means first chan should be 0, secnd 255
			require.Equal(t, 0, sm.GetDMXState().getValue(tt.fields.Universe, tt.fields.StartAddress))
			require.Equal(t, 255, sm.GetDMXState().getValue(tt.fields.Universe, 1+tt.fields.StartAddress))
		})
	}
}

func TestSetLightStateOneStep(t *testing.T) {
	l := DMXLight{
		Name:         "dmxtest",
		StartAddress: 11,
		Universe:     2,
		Profile:      "test_profile",
	}
	s := &config.Server{
		DMXProfiles: config.DMXProfileMap{"test_profile": config.LightProfileDMX{
			Channels: map[string]int{
				"red":   4,
				"green": 5,
				"blue":  11,
			},
		}},
	}
	s.Timings.FadeInterpolationTick = time.Millisecond * 500
	ctx := s.InjectIntoContext(context.Background())
	// m, err := NewManager(ctx, nil)
	m := new(MockManager)
	targetState := TargetState{
		Duration: time.Second,
	}
	redState := State{RGB: color.GetRGBFromString("red")}
	targetState.State = redState

	m.On("GetState", "dmxtest").Return(&State{})
	m.On("SetDMXState",
		mock.Anything,
		dmxOperation{universe: 2, channel: 14, value: 0},
		dmxOperation{universe: 2, channel: 15, value: 0},
		dmxOperation{universe: 2, channel: 21, value: 0},
	).Return(nil)
	m.On("SetState", "dmxtest", State{})

	m.On("SetDMXState",
		mock.Anything,
		dmxOperation{universe: 2, channel: 14, value: 255},
		dmxOperation{universe: 2, channel: 15, value: 0},
		dmxOperation{universe: 2, channel: 21, value: 0},
	).Return(nil)
	m.On("SetState", "dmxtest", redState)
	l.SetState(ctx, m, targetState)

	m.AssertExpectations(t)
}
