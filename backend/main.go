package main

import (
	"fmt"
	"time"

	"github.com/nickysemenza/hyperion/backend/api"
	"github.com/nickysemenza/hyperion/backend/cue"
	"github.com/nickysemenza/hyperion/backend/homekit"
	"github.com/nickysemenza/hyperion/backend/light"
)

func getTempCueStack(CueMaster *cue.Master) cue.Stack {
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
	return mainCueStack
}
func main() {
	fmt.Println("Hello!")

	//read light config
	//TODO: other config like ports and addresses in another file?
	light.ReadLightConfigFromFile("./light/testconfig.json")

	//Set up cue stacks
	cue.CM = cue.Master{}
	CueMaster := &cue.CM
	mainCueStack := getTempCueStack(CueMaster)
	CueMaster.CueStacks = append(CueMaster.CueStacks, mainCueStack)

	//Set up Homekit Server
	go homekit.Start()

	//Set up RPC server
	//go api.ServeRPC(8888)

	//Setup API server
	go api.ServeHTTP()

	//proceess cues forever
	// CueMaster.ProcessForever()
	for {
		time.Sleep(1 * time.Second)
	}
}
