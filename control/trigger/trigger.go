package trigger

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/util/color"
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
func ProcessTriggers(ctx context.Context) {
	c := getTriggerChan()
	for t := range c {

		var newCues []cue.Cue
		log.Printf("new trigger! %v\n", t)
		if t.id == 1 {
			newCues = append(newCues, cue.NewSimple("hue1", color.GetRGBFromString("red")))
			newCues = append(newCues, cue.NewSimple("hue2", color.GetRGBFromString("blue")))
		}
		if t.id == 2 {
			newCues = append(newCues, cue.NewSimple("hue1", color.GetRGBFromString("green")))
		}
		if t.id == 3 {
			newCues = append(newCues, cue.NewSimple("hue1", color.GetRGBFromString("blue")))
		}
		if t.id == 4 {
			newCues = append(newCues, cue.NewSimple("hue1", color.GetRGBFromString("black")))
			newCues = append(newCues, cue.NewSimple("hue2", color.GetRGBFromString("black")))
		}

		for _, x := range newCues {
			stack := cue.GetCueMaster().GetDefaultCueStack()
			stack.EnQueueCue(x)
		}
	}
}
