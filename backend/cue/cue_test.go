package cue

import (
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/backend/light"
	"github.com/stretchr/testify/assert"
)

//Tests getting the Duration of a CueFrame
func TestCueFrameGetDuration(t *testing.T) {
	tt := []struct {
		cf               Frame
		expectedDuration time.Duration
	}{
		{Frame{
			Actions: []FrameAction{
				{NewState: light.State{Duration: time.Millisecond, RGB: light.RGBColor{}}},
			},
		}, time.Millisecond},
		{Frame{
			Actions: []FrameAction{
				{NewState: light.State{Duration: time.Millisecond}},
				{NewState: light.State{Duration: time.Millisecond * 50}},
			},
		}, time.Millisecond * 50},
		{Frame{
			Actions: []FrameAction{},
		}, time.Duration(0)},
	}
	for _, x := range tt {
		if x.cf.GetDuration() != x.expectedDuration {
			t.Errorf("got %d, wanted %d", x.cf.GetDuration(), x.expectedDuration)
		}

		t1 := time.Now()
		x.cf.ProcessFrame()
		t2 := time.Now()
		//5ms of padding/lenience
		assert.WithinDuration(t, t1, t2, x.expectedDuration+(5*time.Millisecond))

		// t.Error(diff)

	}
}
