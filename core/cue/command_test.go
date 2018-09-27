package cue

import (
	"context"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	tt := []struct {
		cmd         string
		expectedCue *Cue
		expectedErr error
	}{
		{"", nil, errorMissingFunction},
		{"a", nil, errorMissingFunction},
		{"foo(bar)", nil, errorUndefinedFunction},
		{"set()", nil, errorWrongPartCount},
		{"set(a:b)", nil, errorWrongPartCount},
		{"set(a:b:)", nil, errorWrongPartCount},
		{"set(light1:green,#0000FF:1000)", nil, errorPartSizeMismatch},
		{"set(light1:green:1 second)", nil, errorInvalidTime},
		{"set(light1:green:1s)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.TargetState{
						Duration: time.Duration(time.Second),
						State:    light.State{RGB: color.RGB{G: 255}},
					}},
			}},
		},
		}, nil},
		{"set(light1:#00FF00:1s|light1:#0000FF:1s)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.TargetState{
						Duration: time.Duration(time.Second),
						State:    light.State{RGB: color.RGB{G: 255}},
					}},
			}},
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.TargetState{
						Duration: time.Duration(time.Second),
						State:    light.State{RGB: color.RGB{B: 255}},
					}},
			}},
		},
		}, nil},
		{"set(light1,light2:#00FF00,#FF0000:1s,2.2s)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				{
					LightName: "light1",
					NewState: light.TargetState{
						Duration: time.Duration(time.Second),
						State:    light.State{RGB: color.RGB{G: 255}},
					},
				},
				{
					LightName: "light2",
					NewState: light.TargetState{
						Duration: time.Duration(time.Millisecond * 2200),
						State:    light.State{RGB: color.RGB{R: 255}},
					},
				},
			}},
		},
		}, nil},
		{"cycle(light1:500ms)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.TargetState{
						Duration: time.Duration(time.Second) / 2,
						State:    light.State{RGB: color.RGB{R: 255}},
					}},
			}},
		},
		}, nil},
	}

	for _, tc := range tt {
		t.Run(tc.cmd, func(t *testing.T) {
			require := require.New(t)
			config := config.Server{}
			ctx := config.InjectIntoContext(context.Background())
			cue, err := NewFromCommand(ctx, tc.cmd)
			if tc.expectedCue == nil {
				require.Nil(cue)
			} else {
				cue.ID = 0
				cue.Status = "" //todo:fix
				for x := range cue.Frames {
					cue.Frames[x].ID = 0
					for y := range cue.Frames[x].Actions {
						cue.Frames[x].Actions[y].ID = 0
						// require.Equal(tc.expectedCue.Frames[x].Actions[y], action)
					}
				}
				require.EqualValues(tc.expectedCue, cue)
			}
			// assert.EqualValues(t, cue, tc.expectedCue)
			require.Equal(err, tc.expectedErr)
		})
	}
}
