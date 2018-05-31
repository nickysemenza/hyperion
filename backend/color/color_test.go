package color

import "testing"

func TestGetRGBComponents(t *testing.T) {
	c := RGB{R: 23, G: 43, B: 0}
	r, g, b := c.AsComponents()
	if r != 23 || g != 43 || b != 0 {
		t.Error("AsComponents is broken!")
	}
}
