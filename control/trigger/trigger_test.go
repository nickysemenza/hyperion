package trigger

import (
	"context"
	"testing"

	"github.com/nickysemenza/hyperion/util/clock"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
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
	m := cue.InitializeMaster(clock.RealClock{})
	stack := m.GetDefaultCueStack()

	//should start with empty cue stack
	assert.Len(t, stack.Cues, 0)
	//first command will add 1
	Action(ctx, "foo", 1)
	assert.Len(t, stack.Cues, 1)
	//bogus command shouldn't touch it
	Action(ctx, "bar", 2)
	assert.Len(t, stack.Cues, 1)
}
