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
