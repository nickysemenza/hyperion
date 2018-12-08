package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	c := RealClock{}
	//not exact, because time changes during execution
	require.WithinDuration(t, c.Now(), time.Now(), time.Millisecond)
}
func TestSmoke(t *testing.T) {
	for _, x := range []Clock{RealClock{}} {
		x.Now()
		x.Sleep(time.Millisecond)
		x.After(time.Millisecond)
	}
	require.True(t, true)
}
