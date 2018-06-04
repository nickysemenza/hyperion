package main

import (
	"github.com/nickysemenza/hyperion/backend/cue"
	"github.com/nickysemenza/hyperion/backend/hyperion"
	"github.com/nickysemenza/hyperion/backend/light"
)

func main() {
	light.ReadLightConfigFromFile("./light/testconfig.json")

	go func() {
		c, _ := cue.BuildCueFromCommand("hue1:#00FF00:1000")
		cs := cue.GetCueMaster().GetDefaultCueStack()
		cs.EnQueueCue(*c)
	}()

	hyperion.RunServer()
}
