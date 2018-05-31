package light

import (
	"testing"
)

func TestDMXAttributeChannels(t *testing.T) {
	tt := []struct {
		profile  dmxProfile
		name     string
		expected int
	}{
		{dmxProfile{Channels: []string{"red", "green"}}, "red", 0},
	}
	for _, tc := range tt {
		res := tc.profile.getChannelIndexForAttribute(tc.name)
		if res != tc.expected {
			t.Errorf("got channel index %d, expected %d", res, tc.expected)
		}
	}
}
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
