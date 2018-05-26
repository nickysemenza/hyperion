package cue

import (
	"testing"
	"time"
)

//Tests getting the Duration of a CueFrame
func TestCueFrameGetDuration(t *testing.T) {
	tt := []struct {
		cf               Frame
		expectedDuration time.Duration
	}{
		{Frame{
			Actions: []FrameAction{
				{Duration: time.Second},
			},
		}, time.Second},
		{Frame{
			Actions: []FrameAction{
				{Duration: time.Second},
				{Duration: time.Second * 3},
			},
		}, time.Second * 3},
		{Frame{
			Actions: []FrameAction{},
		}, time.Duration(0)},
	}
	for _, x := range tt {
		if x.cf.GetDuration() != x.expectedDuration {
			t.Errorf("got %d, wanted %d", x.cf.GetDuration(), x.expectedDuration)
		}
	}
}
