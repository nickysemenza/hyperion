package cue

import (
	"errors"
	"testing"
	"time"

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
		{"set()", nil, errors.New(commandErrorWrongPartCount)},
		{"set(a:b)", nil, errors.New(commandErrorWrongPartCount)},
		{"set(a:b:)", nil, errors.New(commandErrorWrongPartCount)},
		{"set(light1:#00FF00,#0000FF:1000)", nil, errors.New(commandErrorPartSizeMismatch)},
		{"set(light1:#00FF00:1 second)", nil, errors.New(commandErrorInvalidTime)},
		{"set(light1:#00FF00:1000)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.State{
						Duration: time.Duration(time.Second),
						RGB:      color.RGB{G: 255},
					}},
			}},
		},
		}, nil},
		{"set(light1:#00FF00:1000|light1:#0000FF:1000)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.State{
						Duration: time.Duration(time.Second),
						RGB:      color.RGB{G: 255},
					}},
			}},
			{Actions: []FrameAction{
				FrameAction{
					LightName: "light1",
					NewState: light.State{
						Duration: time.Duration(time.Second),
						RGB:      color.RGB{B: 255},
					}},
			}},
		},
		}, nil},
		{"set(light1,light2:#00FF00,#FF0000:1000,2000)", &Cue{Frames: []Frame{
			{Actions: []FrameAction{
				{
					LightName: "light1",
					NewState: light.State{
						Duration: time.Duration(time.Second),
						RGB:      color.RGB{G: 255},
					},
				},
				{
					LightName: "light2",
					NewState: light.State{
						Duration: time.Duration(time.Second * 2),
						RGB:      color.RGB{R: 255},
					},
				},
			}},
		},
		}, nil},
	}

	for _, tc := range tt {
		t.Run(tc.cmd, func(t *testing.T) {
			require := require.New(t)
			cue, err := NewFromCommand(tc.cmd)
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
