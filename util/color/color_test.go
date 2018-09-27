package color

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRGBComponents(t *testing.T) {
	c := RGB{R: 23, G: 43, B: 0}
	r, g, b := c.AsComponents()
	require.Equal(t, r, c.R)
	require.Equal(t, g, c.G)
	require.Equal(t, b, c.B)
}

func TestAsColorfulAndHex(t *testing.T) {
	c1 := RGB{R: 12, G: 255, B: 120}
	require.EqualValues(t, 1, c1.AsColorful().G)
	require.EqualValues(t, "#0cff78", c1.ToHex())
}
func TestIsBlack(t *testing.T) {
	c1 := RGB{R: 12, G: 255, B: 120}
	c2 := RGB{}
	require.False(t, c1.IsBlack())
	require.True(t, c2.IsBlack())
}
func TestAsPB(t *testing.T) {
	c1 := RGB{R: 12, G: 255, B: 120}
	require.EqualValues(t, c1.AsPB().R, c1.R)
	require.EqualValues(t, c1.AsPB().G, c1.G)
	require.EqualValues(t, c1.AsPB().B, c1.B)
}
func TestGetInterpolatedFade(t *testing.T) {
	c1 := RGB{}
	require.Equal(t, 0, c1.GetInterpolatedFade(RGB{R: 200}, 0, 2).R)
	require.Equal(t, 200, c1.GetInterpolatedFade(RGB{R: 200}, 1, 2).R)

	//allow leniency for non-linear fades
	require.Equal(t, 0, c1.GetInterpolatedFade(RGB{R: 33}, 0, 3).R)
	require.InDelta(t, 19, c1.GetInterpolatedFade(RGB{R: 33}, 1, 3).R, 5)
	require.Equal(t, 33, c1.GetInterpolatedFade(RGB{R: 33}, 2, 3).R)

}

func TestGetRGBFromString(t *testing.T) {
	tests := []struct {
		name string
		want RGB
	}{
		//bad inputs
		{"", RGB{}},
		{"foo", RGB{}},
		{"#foo", RGB{}},
		//good inputs
		{"black", RGB{}},
		{"red", RGB{R: 255}},
		{"green", RGB{G: 255}},
		{"blue", RGB{B: 255}},
		{"#FF0000", RGB{R: 255}},
		{"#004A69", RGB{R: 0, G: 73, B: 105}},
		{"white", RGB{R: 255, G: 255, B: 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetRGBFromString(tt.name)
			require.Equal(t, c, tt.want)
		})
	}
}
