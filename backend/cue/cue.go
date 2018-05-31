package cue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/nickysemenza/hyperion/backend/color"
	"github.com/nickysemenza/hyperion/backend/light"
)

var CM Master

//Master is the parent of all CueStacks, is a singleton
type Master struct {
	CueStacks []Stack `json:"cue_stacks"`
	currentID int64
}

//Stack is basically a precedence priority queue (really a CueQueue sigh)
type Stack struct {
	Priority      int64  `json:"priority"`
	Name          string `json:"name"`
	Cues          []Cue  `json:"cues"`
	ProcessedCues []Cue  `json:"processed_cues"`
	Test          *sync.Mutex
	ActiveCue     *Cue `json:"active_cue"`
}

//Cue is a cue.
type Cue struct {
	ID              int64   `json:"id"`
	Frames          []Frame `json:"frames"`
	Name            string  `json:"name"`
	shouldRepeat    bool
	shouldHoldAfter bool //default false, will pause the CueStack after executing this cue, won't move on to next
	waitBefore      time.Duration
	waitAfter       time.Duration
}

//Frame is a single 'animation frame' of a Cue
type Frame struct {
	Actions []FrameAction `json:"actions"`
	ID      int64         `json:"id"`
}

//FrameAction is an action within a Cue(Frame) to be executed simultaneously
type FrameAction struct {
	NewState  light.State `json:"new_state"`
	ID        int64       `json:"id"`
	LightName string      `json:"light_name"`
	//TODO: add `light`
	//TODO: add way to have a noop action (to block aka wait for time)
}

//NewFrameAction creates a new instate with incr ID
func (cm *Master) NewFrameAction(duration time.Duration, color color.RGBColor, lightName string) FrameAction {
	return FrameAction{ID: cm.getNextIDForUse(), LightName: lightName, NewState: light.State{RGB: color, Duration: duration}}
}

//DumpToFile write the CueMaster to a file
func (cm *Master) DumpToFile(fileName string) error {
	jsonData, err := json.MarshalIndent(cm, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, jsonData, 0644)

}

func (cm *Master) getNextIDForUse() int64 {
	//todo:mutex?
	id := cm.currentID
	cm.currentID++
	return id
}

//NewStack makes a new cue stack
func (cm *Master) NewStack(priority int, name string) Stack {
	return Stack{Priority: 2, Name: "main", Test: &sync.Mutex{}}
}

//NewFrame creates a new instate with incr ID
func (cm *Master) NewFrame(actions []FrameAction) Frame {
	return Frame{ID: cm.getNextIDForUse(), Actions: actions}
}

//New creates a new instate with incr ID
func (cm *Master) New(frames []Frame, name string) Cue {
	return Cue{ID: cm.getNextIDForUse(), Frames: frames}
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
		if nextCue := cs.deQueueNextCue(); nextCue != nil {
			cs.ActiveCue = nextCue
			nextCue.ProcessCue()
			cs.ActiveCue = nil
			cs.ProcessedCues = append(cs.ProcessedCues, *nextCue)
		} else {
			// fmt.Println("FINISHED PROCESSING CUESTACK, RESTARTING")
		}
	}
}

func (cs *Stack) deQueueNextCue() *Cue {
	if len(cs.Cues) > 0 {
		cs.Test.Lock()
		x := cs.Cues[0]
		cs.Cues = cs.Cues[1:]
		cs.Test.Unlock()
		return &x
	}
	return nil
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
	//no blocking, so wait until all the child frames have theoretically finished
	time.Sleep(cf.GetDuration())
}

//ProcessFrameAction does job stuff
func (cfa *FrameAction) ProcessFrameAction() {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("[FrameAction #%d] processing @ %d (delta=%s) (color=%v) (light=%s)\n", cfa.ID, now, cfa.NewState.Duration, cfa.NewState.RGB.String(), cfa.LightName)

	if l := light.GetByName(cfa.LightName); l != nil {
		go l.SetState(cfa.NewState)
	} else {
		fmt.Printf("Cannot find light by name: %s\n", cfa.LightName)
	}
	//goroutine doesn't block, so hold until the SetState has (hopefully) finished timing-wise
	time.Sleep(cfa.NewState.Duration)
}
