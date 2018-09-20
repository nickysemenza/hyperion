package trigger

import (
	"context"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/stretchr/testify/assert"
)

func TestTrigger(t *testing.T) {
	//just smoke test for now, make sure channels don't cause deadlock or anything
	go ProcessTriggers(context.Background())
	Action("aa", 1)
	time.Sleep(time.Millisecond * 200) //TODO: make this better
	stack := cue.GetCueMaster().GetDefaultCueStack()

	assert.Len(t, stack.Cues, 2)
	assert.Equal(t, getTriggerChan(), getTriggerChan())
	Action("aa", 2)
	time.Sleep(time.Millisecond * 200)
	assert.Len(t, stack.Cues, 3)
}
