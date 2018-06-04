package light

import (
	"testing"
	"time"
)

func TestGetTransitionTimeAs100msMultiple(t *testing.T) {
	tt := []struct {
		input    time.Duration
		expected uint16
	}{
		{time.Duration(time.Second), 10},
		{time.Duration(time.Millisecond * 200), 2},
		{time.Duration(time.Millisecond * 250), 2},
		{time.Duration(time.Millisecond * 270), 2},
		{time.Duration(0), 0},
	}
	for _, x := range tt {
		res := getTransitionTimeAs100msMultiple(x.input)
		if res != x.expected {
			t.Errorf("got %d, expected %d", res, x.expected)
		}
	}
}
