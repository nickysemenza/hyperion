package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/hyperion/backend/api"
	"github.com/nickysemenza/hyperion/backend/color"
	"github.com/nickysemenza/hyperion/backend/cue"
	"github.com/nickysemenza/hyperion/backend/homekit"
	"github.com/nickysemenza/hyperion/backend/light"
)

func getTempCueStack(CueMaster *cue.Master) cue.Stack {
	mainCueStack := CueMaster.NewStack(2, "main")
	// for x := 1; x <= 2; x++ {
	// 	a := CueMaster.New([]cue.Frame{
	// 		CueMaster.NewFrame([]cue.FrameAction{
	// 			CueMaster.NewFrameAction(time.Millisecond*1500, color.RGB{R: 255}, "hue1"),
	// 			CueMaster.NewFrameAction(0, color.RGB{R: 255}, "hue2"),
	// 		}),
	// 		CueMaster.NewFrame([]cue.FrameAction{
	// 			CueMaster.NewFrameAction(time.Second*time.Duration(x), color.RGB{G: 255}, "hue1"),
	// 			CueMaster.NewFrameAction(0, color.RGB{B: 255}, "hue2"),
	// 		}),
	// 		CueMaster.NewFrame([]cue.FrameAction{
	// 			CueMaster.NewFrameAction(0, color.RGB{B: 255}, "hue1"),
	// 			CueMaster.NewFrameAction(0, color.RGB{R: 255}, "hue2"),
	// 		}),
	// 		CueMaster.NewFrame([]cue.FrameAction{
	// 			CueMaster.NewFrameAction(time.Second*2, color.RGB{B: 255}, "hue1"),
	// 			CueMaster.NewFrameAction(time.Second*2, color.RGB{B: 255}, "hue2"),
	// 		}),
	// 	}, fmt.Sprintf("Cue #%d", x))
	// 	mainCueStack.Cues = append(mainCueStack.Cues, a)
	// }
	a := CueMaster.New([]cue.Frame{
		CueMaster.NewFrame([]cue.FrameAction{
			CueMaster.NewFrameAction(time.Millisecond*2500, color.RGB{0, 255, 100}, "par1"),
			CueMaster.NewFrameAction(time.Millisecond*2500, color.RGB{0, 255, 100}, "par2"),
		}),
		CueMaster.NewFrame([]cue.FrameAction{
			CueMaster.NewFrameAction(time.Millisecond*8500, color.RGB{255, 111, 37}, "par1"),
		}),
		CueMaster.NewFrame([]cue.FrameAction{
			CueMaster.NewFrameAction(time.Millisecond*8500, color.RGB{1, 1, 255}, "par1"),
		}),
	}, "Cue 1")
	mainCueStack.Cues = append(mainCueStack.Cues, a)
	return mainCueStack
}
func main() {
	fmt.Println("Hello!")

	//read light config
	//TODO: other config like ports and addresses in another file?
	light.ReadLightConfigFromFile("./light/testconfig.json")

	spew.Dump(light.Config)

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
	CueMaster.ProcessForever()

	go light.SendDMXValuesToOLA()

	//handle CTRL+C
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Println("Shutdown hyperion ...")
		os.Exit(0)
	}()

	//keep going
	for {
		time.Sleep(1 * time.Second)
	}
}
