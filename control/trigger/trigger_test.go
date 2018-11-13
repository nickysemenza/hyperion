package trigger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nickysemenza/hyperion/util/clock"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/stretchr/testify/assert"
)

func TestTrigger(t *testing.T) {
	//just smoke test for now, make sure channels don't cause deadlock or anything
	config := config.Server{
		Triggers: []config.Trigger{
			{ID: 1, Source: "foo", Command: "set(hue2:blue:1s)"},
			{ID: 2, Source: "bar", Command: "bogus"},
		},
	}
	ctx := config.InjectIntoContext(context.Background())
	m := cue.InitializeMaster(clock.RealClock{}, &light.Manager{})
	stack := m.GetDefaultCueStack()

	//should start with empty cue stack
	assert.Len(t, stack.Cues, 0)
	//first command will add 1
	Action(ctx, "foo", 1, m)
	assert.Len(t, stack.Cues, 1)
	//bogus command shouldn't touch it
	Action(ctx, "bar", 2, m)
	assert.Len(t, stack.Cues, 1)

	m2 := new(cue.MockMaster)
	cs := cue.Stack{Name: "foo"}
	m2.On("GetDefaultCueStack").Return(&cs)
	newCue, err := cue.NewFromCommand(ctx, "set(hue2:blue:1s)")
	newCue.Source = cue.Source{
		Input: "trigger",
		Type:  "command",
		Meta:  "trigger=foo:1",
	}
	require.NoError(t, err)
	m2.On("EnQueueCue", *newCue, &cs).Return(newCue)

	Action(ctx, "foo", 1, m2)
	m2.AssertExpectations(t)
}
