package color

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//RGB holds RGB values (0-255, although they can be negative in the case of a delta)
type RGB struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type colorNum int

//different hardcoded colors
const (
	Red colorNum = iota
	Green
	Blue
	White
)

//FromString yields an RGB object based on a const
func FromString(num colorNum) RGB {
	switch num {
	case Red:
		return RGB{R: 255}
	case Green:
		return RGB{G: 255}
	case Blue:
		return RGB{B: 255}
	case White:
		return RGB{R: 255, G: 255, B: 255}
	default:
		return RGB{}
	}
}

//GetInterpolatedFade returns fade from one color to another.
func (c *RGB) GetInterpolatedFade(target RGB, step, numSteps int) colorful.Color {
	c1 := c.AsColorful()
	c2 := target.AsColorful()

	return c1.BlendHcl(c2, float64(step)/float64(numSteps-1)).Clamped()
}

//AsColorful turns a colorful.Color into an RGB
func (c *RGB) AsColorful() colorful.Color {
	return colorful.Color{R: float64(c.R) / 255, G: float64(c.G) / 255, B: float64(c.B) / 255}
}

//GetRGBFromColorful turns a colorful.Color struct into an RGB one
func GetRGBFromColorful(c colorful.Color) RGB {
	return RGB{
		R: int(c.R * 255),
		G: int(c.G * 255),
		B: int(c.B * 255),
	}
}

//GetRGBFromHex turns a hex string (#00FF00) into an RGB color
func GetRGBFromHex(hex string) RGB {
	c, err := colorful.Hex(hex)
	if err != nil {
		log.Println(err)
	}
	return GetRGBFromColorful(c)
}

//AsComponents returns the seperate r, g, b
func (c *RGB) AsComponents() (int, int, int) {
	return c.R, c.G, c.B
}

//String returns a ANSI-color formatted r/g/b string
func (c *RGB) String() string {

	red := color.New(color.BgRed).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	blue := color.New(color.BgBlue).SprintFunc()
	return fmt.Sprintf("%s %s %s", red(c.R), green(c.G), blue(c.B))
}

//GetXyy returns the RGB color in xyy color space
func (c *RGB) GetXyy() (x, y, Yout float64) {
	cc := colorful.Color{
		R: float64(c.R / 255),
		G: float64(c.G / 255),
		B: float64(c.B / 255),
	}
	//OLD:  x, y, _ := cc.Xyz()
	return colorful.XyzToXyy(colorful.LinearRgbToXyz(cc.LinearRgb()))
}
