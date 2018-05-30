package light

import "testing"

func TestRGBColorDelta(t *testing.T) {
	tt := []struct {
		name string
		from RGBColor
		to   RGBColor
		res  RGBColor
	}{
		{
			"all@255 => all@0 == delta all@-255",
			RGBColor{255, 255, 255},
			RGBColor{0, 0, 0},
			RGBColor{-255, -255, -255},
		},
		{
			"all@0 => all@255 == delta all@255",
			RGBColor{0, 0, 0},
			RGBColor{255, 255, 255},
			RGBColor{255, 255, 255},
		},
		{
			"same should yield zero delta",
			RGBColor{25, 184, 4},
			RGBColor{25, 184, 4},
			RGBColor{0, 0, 0},
		},
		{
			"mixed direction",
			RGBColor{25, 184, 4},
			RGBColor{25, 4, 29},
			RGBColor{0, -180, 25},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			delta := tc.from.DeltaTo(tc.to)
			if delta != tc.res {
				t.Errorf("expected delta to be %s, got %s", tc.res.String(), delta.String())
			}
		})

	}
}
