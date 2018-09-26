package cue

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
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
				{NewState: light.TargetState{Duration: time.Millisecond, State: light.State{RGB: color.RGB{}}}},
			},
		}, time.Millisecond},
		{Frame{
			Actions: []FrameAction{
				{NewState: light.TargetState{Duration: time.Millisecond}},
				{NewState: light.TargetState{Duration: time.Millisecond * 50}},
			},
		}, time.Millisecond * 50},
		{Frame{
			Actions: []FrameAction{},
		}, time.Duration(0)},
	}
	for _, x := range tt {
		require.Equal(t, x.expectedDuration, x.cf.GetDuration())
		//test with real timings
		t1 := time.Now()
		x.cf.ProcessFrame(context.Background())
		t2 := time.Now()
		//7ms of padding/lenience (CI is slow)
		require.WithinDuration(t, t1, t2, x.expectedDuration+(7*time.Millisecond))
	}
}

func TestCueDurationHelpers(t *testing.T) {
	tests := []struct {
		c                     Cue
		expectedDuration      time.Duration
		expectedDurationDrift time.Duration
	}{
		{Cue{
			RealDuration: time.Millisecond * 3,
			Status:       statusProcessed,
		}, 0, time.Millisecond * 3},
		{Cue{
			RealDuration: time.Millisecond * 3,
			Status:       statusActive,
		}, 0, 0},
		{Cue{
			Status:       statusProcessed,
			RealDuration: time.Millisecond * 25,
			Frames: []Frame{
				{Actions: []FrameAction{
					{NewState: light.TargetState{Duration: time.Millisecond * 7}},
					{NewState: light.TargetState{Duration: time.Millisecond * 12}},
				}},
				{Actions: []FrameAction{
					{NewState: light.TargetState{Duration: time.Millisecond * 8}},
					{NewState: light.TargetState{Duration: time.Millisecond * 11}},
				}},
			},
		}, time.Millisecond * 23, time.Millisecond * 2},
	}

	for _, tt := range tests {
		require := require.New(t)
		cue := &tt.c
		require.Equal(tt.expectedDuration, cue.GetDuration())
		require.Equal(tt.expectedDurationDrift, cue.getDurationDrift())

		if cue.Status != statusActive {
			require.Zero(cue.getElapsedTime())
		} else {
			cue.StartedAt = time.Now()
			time.Sleep(time.Microsecond)
			require.NotZero(cue.getElapsedTime())
		}

		t1 := time.Now()
		cue.ProcessCue(context.Background())
		t2 := time.Now()
		//TODO: move status switching form ProcessStack to ProcessCue
		// require.Equal(statusProcessed, cue.Status)

		//7ms of padding/lenience (CI is slow)
		require.WithinDuration(t1, t2, tt.expectedDuration+(7*time.Millisecond))
	}
}

func TestCueQueueing(t *testing.T) {
	cs := Stack{}
	assert.Nil(t, cs.deQueueNextCue(), "deque on empty should return nil")

	c1 := cs.EnQueueCue(Cue{Name: "c1"})
	c2 := cs.EnQueueCue(Cue{Name: "c2"})
	require.NotEqual(t, c1.ID, c2.ID)

	assert.Equal(t, len(cs.Cues), 2)
	pop := cs.deQueueNextCue()
	assert.NotNil(t, pop)

	assert.Equal(t, pop.Name, "c1", "queue should be FIFO")

	cs.EnQueueCue(Cue{Name: "c3"})
	assert.Equal(t, cs.deQueueNextCue().Name, "c2", "queue should be FIFO")
	assert.Equal(t, cs.deQueueNextCue().Name, "c3", "queue should be FIFO")

	assert.Nil(t, cs.deQueueNextCue())
}

func TestCueMarshalling(t *testing.T) {
	require := require.New(t)

	//FrameAction
	cfa := FrameAction{NewState: light.TargetState{Duration: time.Millisecond * 7}, ID: 1}
	b, err := cfa.MarshalJSON()
	require.NoError(err)
	json := fmt.Sprintf("%s", b)
	require.Contains(json, `"action_duration_ms":7`)
	require.Contains(json, `"id":1`)

	//Frame
	cf := Frame{Actions: []FrameAction{
		FrameAction{NewState: light.TargetState{Duration: time.Millisecond * 8}, ID: 2},
		FrameAction{NewState: light.TargetState{Duration: time.Millisecond * 9}, ID: 3},
	}}
	b, err = cf.MarshalJSON()
	require.NoError(err)
	json = fmt.Sprintf("%s", b)
	require.Contains(json, `"expected_duration_ms":9`)
	require.Contains(json, `"id":2`)

	//Cue
	c := Cue{
		Status:       statusProcessed,
		RealDuration: time.Millisecond * 25,
		Frames: []Frame{
			{Actions: []FrameAction{
				{NewState: light.TargetState{Duration: time.Millisecond * 7}},
				{NewState: light.TargetState{Duration: time.Millisecond * 12}},
			}},
			{Actions: []FrameAction{
				{NewState: light.TargetState{Duration: time.Millisecond * 8}},
				{NewState: light.TargetState{Duration: time.Millisecond * 11}},
			}},
		},
	}
	b, err = c.MarshalJSON()
	require.NoError(err)
	json = fmt.Sprintf("%s", b)
	require.Contains(json, `"duration_drift_ms":2`)
}

func BenchmarkCueFrameProcessing(b *testing.B) {
	actions := []FrameAction{}
	for i := 0; i < b.N; i++ {
		actions = append(actions, FrameAction{NewState: light.TargetState{Duration: 0, State: light.State{RGB: color.RGB{}}}})
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
