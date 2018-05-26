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

	spew.Dump(light.ReadLightConfigFromFile("./light/testconfig.json"))

	// os.Exit(0)

	cm := &cue.Master{}

	mainCueStack := cue.Stack{Priority: 2, Name: "main"}
	for x := 1; x <= 2; x++ {
		a := cm.New([]cue.Frame{
			cm.NewFrame([]cue.FrameAction{
				cm.NewFrameAction(time.Millisecond*1500, cue.RGBColor{}),
				cm.NewFrameAction(time.Second*time.Duration(x), cue.RGBColor{}),
			})}, fmt.Sprintf("Cue #%d", x))
		mainCueStack.Cues = append(mainCueStack.Cues, a)
	}

	// secondaryCuestack := CueStack{}
	// copier.Copy(&secondaryCuestack, &mainCueStack)
	// secondaryCuestack.Name = "secondary"

	cm.CueStacks = append(cm.CueStacks, mainCueStack)
	spew.Dump(cm)
	cm.ProcessForever()
	fmt.Println("faaa")

	for {
		time.Sleep(1 * time.Second)
	}
}
