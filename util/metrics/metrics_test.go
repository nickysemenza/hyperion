package metrics

import (
	"testing"
	"time"
)

func TestMetricsSimple(t *testing.T) {
	Register()                                            //could panic
	SetGagueWithNsFromTime(time.Now(), ResponseTimeNsHue) //no idea how to test this...
}
