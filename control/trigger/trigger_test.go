package trigger

import (
	"context"
	"testing"

	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/stretchr/testify/assert"
)

func TestTrigger(t *testing.T) {
	//just smoke test for now, make sure channels don't cause deadlock or anything
	ctx := context.Background()
	stack := cue.GetCueMaster().GetDefaultCueStack()

	Action(ctx, "aa", 1)

	assert.Len(t, stack.Cues, 2)
	Action(ctx, "aa", 2)
	assert.Len(t, stack.Cues, 3)
}
