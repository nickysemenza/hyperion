package color

import (
	"fmt"

	"github.com/fatih/color"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//RGBColor holds RGB values (0-255, although they can be negative in the case of a delta)
type RGBColor struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

func (c *RGBColor) AsColorful() colorful.Color {
	return colorful.Color{R: float64(c.R) / 255, G: float64(c.G) / 255, B: float64(c.B) / 255}
}

func GetRGBFromColorful(c colorful.Color) RGBColor {
	return RGBColor{
		R: int(c.R * 255),
		G: int(c.G * 255),
		B: int(c.B * 255),
	}
}

//AsComponents returns the seperate r, g, b
func (c *RGBColor) AsComponents() (int, int, int) {
	return c.R, c.G, c.B
}

//String returns a ANSI-color formatted r/g/b string
func (c *RGBColor) String() string {

	red := color.New(color.BgRed).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	blue := color.New(color.BgBlue).SprintFunc()
	return fmt.Sprintf("%s %s %s", red(c.R), green(c.G), blue(c.B))
}

//GetXyy returns the RGB color in xyy color space
func (c *RGBColor) GetXyy() (x, y, Yout float64) {
	cc := colorful.Color{
		R: float64(c.R / 255),
		G: float64(c.G / 255),
		B: float64(c.B / 255),
	}
	//OLD:  x, y, _ := cc.Xyz()
	return colorful.XyzToXyy(colorful.LinearRgbToXyz(cc.LinearRgb()))
}
