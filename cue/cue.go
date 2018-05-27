package cue

import (
	"fmt"
	"log"
	"time"

	"github.com/nickysemenza/hyperion/light"
)

//Master is the parent of all CueStacks, is a singleton
type Master struct {
	CueStacks  []Stack
	CurrentIDs struct {
		CueStack       int64
		Cue            int64
		CueFrame       int64
		CueFrameAction int64
	}
}

//Stack is basically a precedence priority queue (really a CueQueue sigh)
type Stack struct {
	Priority int64
	Name     string
	Cues     []Cue
}

//Cue is a cue.
type Cue struct {
	ID              int64
	Frames          []Frame
	Name            string
	shouldRepeat    bool
	shouldHoldAfter bool //default false, will pause the CueStack after executing this cue, won't move on to next
	waitBefore      time.Duration
	waitAfter       time.Duration
}

//Frame is a single 'animation frame' of a Cue
type Frame struct {
	Actions []FrameAction
	ID      int64
}

//FrameAction is an action within a Cue(Frame) to be executed simultaneously
type FrameAction struct {
	NewState  light.State
	ID        int64
	LightName string
	//TODO: add `light`
	//TODO: add way to have a noop action (to block aka wait for time)
}

//NewFrameAction creates a new instate with incr ID
func (cm *Master) NewFrameAction(duration time.Duration, color light.RGBColor, lightName string) FrameAction {
	id := cm.CurrentIDs.CueFrameAction
	cm.CurrentIDs.CueFrameAction++
	return FrameAction{ID: id, LightName: lightName, NewState: light.State{RGB: color, Duration: duration}}
}

//NewFrame creates a new instate with incr ID
func (cm *Master) NewFrame(actions []FrameAction) Frame {
	id := cm.CurrentIDs.CueFrame
	cm.CurrentIDs.CueFrame++
	return Frame{ID: id, Actions: actions}
}

//New creates a new instate with incr ID
func (cm *Master) New(frames []Frame, name string) Cue {
	id := cm.CurrentIDs.Cue
	cm.CurrentIDs.Cue++
	return Cue{ID: id, Frames: frames}
}

//ProcessForever runs all the cuestacks
func (cm *Master) ProcessForever() {
	for x := range cm.CueStacks {
		go cm.CueStacks[x].ProcessStack()
	}
}

//ProcessStack processes cues
func (cs *Stack) ProcessStack() {
	log.Printf("[CueStack: %s]\n", cs.Name)
	for {
		for _, eachCue := range cs.Cues {
			eachCue.ProcessCue()
		}
		fmt.Println("FINISHED PROCESSING CUESTACK, RESTARTING")
	}
}

//ProcessCue processes cue
func (c *Cue) ProcessCue() {
	log.Printf("[ProcessCue #%d]\n", c.ID)
	for _, eachFrame := range c.Frames {
		eachFrame.ProcessFrame()
	}
}

//GetDuration returns the longest lasting Action within a CueFrame
func (cf *Frame) GetDuration() time.Duration {
	longest := time.Duration(0)
	for _, action := range cf.Actions {
		if d := action.NewState.Duration; d > longest {
			longest = d
		}
	}
	return longest
}

//ProcessFrame processes the cueframe
func (cf *Frame) ProcessFrame() {
	log.Printf("[CF #%d] Has %d Actions, will take %s\n", cf.ID, len(cf.Actions), cf.GetDuration())
	// fmt.Println(cf.Actions)
	for x := range cf.Actions {
		go cf.Actions[x].ProcessFrameAction()
	}
	time.Sleep(cf.GetDuration())
}

//ProcessFrameAction does job stuff
func (cfa *FrameAction) ProcessFrameAction() {
	//TODO: send dmx, call hue func, etc
	now := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("[FrameAction #2] processing @ %d (delta=%s) (color=%v) (light=%s)", now, cfa.NewState.Duration, cfa.NewState.RGB.FancyString(), cfa.LightName)

	//TODO: switch based on light type/ID
	// br := hue.Bridge{
	// 	Hostname: "10.0.1.55",
	// 	Username: "alW0LsA1mnXB28T4txGs01BeHi1WBr661VZ1eqEF",
	// }

	// LightConfig

	l := light.GetLightByName(cfa.LightName)
	if l != nil {
		go l.SetColor(cfa.NewState.RGB)
	} else {
		fmt.Printf("Cannot find lighty by name: %s", cfa.LightName)
	}

	// go br.SetColor(2, cfa.Color, cfa.Duration)
	// go br.SetColor(1, cfa.Color, cfa.Duration)

	time.Sleep(cfa.NewState.Duration)

}
