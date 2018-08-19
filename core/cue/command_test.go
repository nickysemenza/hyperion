package cue

import (
	"errors"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	tt := []struct {
		cmd         string
		expectedCue *Cue
		expectedErr error
	}{
		{"", nil, errors.New(commandErrorWrongPartCount)},
		{"a:b", nil, errors.New(commandErrorWrongPartCount)},
		{"a:b:", nil, errors.New(commandErrorWrongPartCount)},
		{"light1:#00FF00,#0000FF:1000", nil, errors.New(commandErrorPartSizeMismatch)},
		{"light1:#00FF00:1 second", nil, errors.New(commandErrorInvalidTime)},
		{"light1:#00FF00:1000", &Cue{Frames: []Frame{
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
		{"light1,light2:#00FF00,#FF0000:1000,2000", &Cue{Frames: []Frame{
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
			cue, err := buildCueFromCommand(tc.cmd)
			assert.Equal(t, cue, tc.expectedCue)
			assert.Equal(t, err, tc.expectedErr)
		})
	}
}
