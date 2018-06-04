package trigger

import (
	"log"
	"sync"

	"github.com/nickysemenza/hyperion/color"
	"github.com/nickysemenza/hyperion/cue"
)

type trigger struct {
	source string
	id     int
}

type chanOfTriggers chan trigger

var (
	triggers chanOfTriggers
	once     sync.Once
)

func getTriggerChan() chanOfTriggers {
	once.Do(func() {
		triggers = make(chanOfTriggers, 100)
	})
	return triggers
}

//Action is called when an trigger needs to be fired
func Action(source string, id int) {
	c := getTriggerChan()
	c <- trigger{source, id}
}

//ProcessTriggers is a worker that processes triggers
func ProcessTriggers() {
	c := getTriggerChan()
	for t := range c {
		log.Printf("new trigger! %v\n", t)
		var newCue cue.Cue
		sendCue := false
		if t.id == 1 {
			newCue = cue.NewSimple("hue1", color.FromString(color.Red))
			sendCue = true
		}
		if t.id == 2 {
			newCue = cue.NewSimple("hue1", color.FromString(color.Green))
			sendCue = true
		}
		if t.id == 3 {
			newCue = cue.NewSimple("hue1", color.FromString(color.Blue))
			sendCue = true
		}

		if sendCue {
			stack := cue.GetCueMaster().GetDefaultCueStack()
			stack.EnQueueCue(newCue)
		}
	}
}
