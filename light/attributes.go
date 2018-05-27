package light

import (
	"fmt"

	"github.com/fatih/color"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//RGBColor holds RGB values (0-255)
type RGBColor struct {
	R int
	G int
	B int
}

func (c *RGBColor) FancyString() string {

	red := color.New(color.BgRed).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	blue := color.New(color.BgBlue).SprintFunc()
	return fmt.Sprintf("%s %s %s", red(c.R), green(c.G), blue(c.B))
}

func (c *RGBColor) GetXyy() (x, y, Yout float64) {
	cc := colorful.Color{
		R: float64(c.R / 255),
		G: float64(c.G / 255),
		B: float64(c.B / 255),
	}
	//OLD:  x, y, _ := cc.Xyz()
	return colorful.XyzToXyy(colorful.LinearRgbToXyz(cc.LinearRgb()))
}
