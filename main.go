package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/hyperion/cue"
	"github.com/nickysemenza/hyperion/light"
)

func main() {

	fmt.Println("Hello!")

	light.ReadLightConfigFromFile("./light/testconfig.json")
	spew.Dump(light.Config)

	CueMaster := &cue.Master{}

	mainCueStack := cue.Stack{Priority: 2, Name: "main"}
	for x := 1; x <= 2; x++ {
		a := CueMaster.New([]cue.Frame{
			CueMaster.NewFrame([]cue.FrameAction{
				CueMaster.NewFrameAction(time.Millisecond*1500, light.RGBColor{R: 255}, "hue1"),
				CueMaster.NewFrameAction(0, light.RGBColor{R: 255}, "hue2"),
			}),
			CueMaster.NewFrame([]cue.FrameAction{
				CueMaster.NewFrameAction(time.Second*time.Duration(x), light.RGBColor{G: 255}, "hue1"),
				CueMaster.NewFrameAction(0, light.RGBColor{B: 255}, "hue2"),
			}),
			CueMaster.NewFrame([]cue.FrameAction{
				CueMaster.NewFrameAction(0, light.RGBColor{B: 255}, "hue1"),
				CueMaster.NewFrameAction(0, light.RGBColor{R: 255}, "hue2"),
			}),
			CueMaster.NewFrame([]cue.FrameAction{
				CueMaster.NewFrameAction(time.Second*2, light.RGBColor{B: 255}, "hue1"),
				CueMaster.NewFrameAction(0, light.RGBColor{B: 255}, "hue2"),
			}),
		}, fmt.Sprintf("Cue #%d", x))
		mainCueStack.Cues = append(mainCueStack.Cues, a)
	}

	// secondaryCuestack := CueStack{}
	// copier.Copy(&secondaryCuestack, &mainCueStack)
	// secondaryCuestack.Name = "secondary"

	CueMaster.CueStacks = append(CueMaster.CueStacks, mainCueStack)
	spew.Dump(CueMaster)
	CueMaster.ProcessForever()
	fmt.Println("faaa")

	for {
		time.Sleep(1 * time.Second)
	}
}
