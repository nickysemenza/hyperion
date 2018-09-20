package trigger

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/nickysemenza/hyperion/core/cue"
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
			c1, _ := cue.NewFromCommand("set(hue1:red:1000)")
			c2, _ := cue.NewFromCommand("set(hue2:blue:1000)")
			newCues = append(newCues, *c1, *c2)
		}
		if t.id == 2 {
			c1, _ := cue.NewFromCommand("set(hue1:green:1000)")
			newCues = append(newCues, *c1)
		}
		if t.id == 3 {
			c1, _ := cue.NewFromCommand("set(hue1:blue:1000)")
			newCues = append(newCues, *c1)
		}
		if t.id == 4 {
			c1, _ := cue.NewFromCommand("set(hue1:black:1000)")
			c2, _ := cue.NewFromCommand("set(hue2:black:1000)")
			newCues = append(newCues, *c1, *c2)
		}

		for _, x := range newCues {
			stack := cue.GetCueMaster().GetDefaultCueStack()
			stack.EnQueueCue(x)
		}
	}
}
