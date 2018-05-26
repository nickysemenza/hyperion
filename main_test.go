package main

import (
	"testing"
	"time"
)

//Tests getting the Duration of a CueFrame
func TestCueFrameGetDuration(t *testing.T) {
	tt := []struct {
		cf               CueFrame
		expectedDuration time.Duration
	}{
		{CueFrame{
			Actions: []CueFrameAction{
				{Duration: time.Second},
			},
		}, time.Second},
		{CueFrame{
			Actions: []CueFrameAction{
				{Duration: time.Second},
				{Duration: time.Second * 3},
			},
		}, time.Second * 3},
		{CueFrame{
			Actions: []CueFrameAction{},
		}, time.Duration(0)},
	}
	for _, x := range tt {
		if x.cf.GetDuration() != x.expectedDuration {
			t.Errorf("got %d, wanted %d", x.cf.GetDuration(), x.expectedDuration)
		}
	}
}
