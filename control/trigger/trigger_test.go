package trigger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrigger(t *testing.T) {
	//just smoke test for now, make sure channels don't cause deadlock or anything
	Action("aa", 1)
	assert.Equal(t, getTriggerChan(), getTriggerChan())
	Action("aa", 2)
	go ProcessTriggers()
}
