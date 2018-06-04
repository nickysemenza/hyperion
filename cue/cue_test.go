package cue

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/color"
	"github.com/nickysemenza/hyperion/light"
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
				{NewState: light.State{Duration: time.Millisecond, RGB: color.RGB{}}},
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
		x.cf.ProcessFrame(context.Background())
		t2 := time.Now()
		//7ms of padding/lenience (travis is slow)
		assert.WithinDuration(t, t1, t2, x.expectedDuration+(7*time.Millisecond))

		// t.Error(diff)

	}
}

func TestCueQueueing(t *testing.T) {
	cs := Stack{}
	assert.Nil(t, cs.deQueueNextCue(), "deque on empty should return nil")

	cs.EnQueueCue(Cue{Name: "c1"})
	cs.EnQueueCue(Cue{Name: "c2"})

	assert.Equal(t, len(cs.Cues), 2)
	pop := cs.deQueueNextCue()
	assert.NotNil(t, pop)

	assert.Equal(t, pop.Name, "c1", "queue should be FIFO")

	cs.EnQueueCue(Cue{Name: "c3"})
	assert.Equal(t, cs.deQueueNextCue().Name, "c2", "queue should be FIFO")
	assert.Equal(t, cs.deQueueNextCue().Name, "c3", "queue should be FIFO")

	assert.Nil(t, cs.deQueueNextCue())
}

func BenchmarkCueFrameProcessing(b *testing.B) {
	actions := []FrameAction{}
	for i := 0; i < b.N; i++ {
		actions = append(actions, FrameAction{NewState: light.State{Duration: 0, RGB: color.RGB{}}})
	}
	frame := Frame{Actions: actions}
	frame.ProcessFrame(context.Background())
}

func TestAddingIDsToUnmarshalledCue(t *testing.T) {
	data := `{
		"frames": [
		  {
			"actions": [
			  {
				"new_state": {
				  "rgb": {
					"r": 65,
					"g": 0,
					"b": 120
				  },
				  "duration": 1500000000
				},
				"light_name": "hue1"
			  }
			]
		  }
		],
		"name": ""
	  }`
	cue := Cue{}
	json.Unmarshal([]byte(data), &cue)

	assert.Zero(t, cue.ID)

	cue.AddIDsRecursively()

	assert.NotZero(t, cue.ID)
	assert.NotZero(t, cue.Frames[0].ID)
	assert.NotZero(t, cue.Frames[0].Actions[0].ID)
	assert.NotEqual(t, cue.ID, cue.Frames[0].ID)
}
