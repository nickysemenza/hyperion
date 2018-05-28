package light

import (
	"testing"
)

func TestDMX(t *testing.T) {
	s1 := getDMXStateInstance()
	s1.setDMXValue(2, 22, 40)

	s2 := getDMXStateInstance()
	if s2.universes[2][22-1] != 40 {
		t.Error("didn't set DMX state instance properly")
	}

	if err := s2.setDMXValue(2, 0, 2); err == nil {
		t.Error("should not allow channel 0")
	}

	if s1 != s2 {
		t.Error("should be singleton!")
	}
}
