package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/hyperion/light"
)

//CueMaster is the parent of all CueStacks, is a singleton
type CueMaster struct {
	CueStacks  []CueStack
	CurrentIDs struct {
		CueStack       int64
		Cue            int64
		CueFrame       int64
		CueFrameAction int64
	}
}

//CueStack is basically a precedence priority queue (really a CueQueue sigh)
type CueStack struct {
	Priority int64
	Name     string
	Cues     []Cue
}

//Cue is a cue.
type Cue struct {
	ID              int64
	Frames          []CueFrame
	Name            string
	shouldRepeat    bool
	shouldHoldAfter bool //default false, will pause the CueStack after executing this cue, won't move on to next
	waitBefore      time.Duration
	waitAfter       time.Duration
}

//CueFrame is a single 'animation frame' of a Cue
type CueFrame struct {
	Actions []CueFrameAction
	ID      int64
}

//CueFrameAction is an action within a Cue(Frame) to be executed simultaneously
type CueFrameAction struct {
	Duration time.Duration
	Color    RGBColor
	ID       int64
}

//RGBColor holds RGB values (0-255)
type RGBColor struct {
	R int
	G int
	B int
}

func (cm *CueMaster) NewCueFrameAction(duration time.Duration, color RGBColor) CueFrameAction {
	id := cm.CurrentIDs.CueFrameAction
	cm.CurrentIDs.CueFrameAction++
	return CueFrameAction{ID: id, Duration: duration, Color: color}
}
func (cm *CueMaster) NewCueFrame(actions []CueFrameAction) CueFrame {
	id := cm.CurrentIDs.CueFrame
	cm.CurrentIDs.CueFrame++
	return CueFrame{ID: id, Actions: actions}
}
func (cm *CueMaster) NewCue(frames []CueFrame, name string) Cue {
	id := cm.CurrentIDs.Cue
	cm.CurrentIDs.Cue++
	return Cue{ID: id, Frames: frames}
}

//ProcessForever runs all the cuestacks
func (cm *CueMaster) ProcessForever() {
	for x := range cm.CueStacks {
		go cm.CueStacks[x].ProcessCueStack()
	}
}

//ProcessCueStack processes cues
func (cs *CueStack) ProcessCueStack() {
	fmt.Printf("Processing CueStack: %s\n", cs.Name)
	for {
		for _, eachCue := range cs.Cues {
			eachCue.ProcessCue(cs)
		}
		fmt.Println("FINISHED PROCESSING CUESTACK, RESTARTING")
	}
}

//ProcessCue processes cue
func (c *Cue) ProcessCue(parentCueStack *CueStack) {
	fmt.Printf("[%s][%s] Process()\n", parentCueStack.Name, c.Name)
	for _, eachFrame := range c.Frames {
		eachFrame.ProcessCueFrame(parentCueStack, c)
	}
}

//GetDuration returns the longest lasting Action within a CueFrame
func (cf *CueFrame) GetDuration() time.Duration {
	longest := time.Duration(0)
	for _, action := range cf.Actions {
		if action.Duration > longest {
			longest = action.Duration
		}
	}
	return longest
}

//ProcessCueFrame processes the cueframe
func (cf *CueFrame) ProcessCueFrame(parentCueStack *CueStack, parentCue *Cue) {
	fmt.Printf("[%s][%s][CF] Has %d Actions, will take %s\n", parentCueStack.Name, parentCue.Name, len(cf.Actions), cf.GetDuration())
	// fmt.Println(cf.Actions)
	for x := range cf.Actions {
		go cf.Actions[x].ProcessCueFrameAction(parentCueStack, parentCue)
	}
	time.Sleep(cf.GetDuration())
}

//ProcessCueFrameAction does job stuff
func (cfa *CueFrameAction) ProcessCueFrameAction(parentCueStack *CueStack, parentCue *Cue) {
	fmt.Printf("[%s][%s][CF][CFA: %d], %v\n", parentCueStack.Name, parentCue.Name, cfa.ID, cfa.Duration)
	time.Sleep(cfa.Duration)
	fmt.Printf("[%s][%s][CF][CFA: %d], done\n", parentCueStack.Name, parentCue.Name, cfa.ID)
}

func main() {

	fmt.Println("Hello!")

	spew.Dump(light.ReadLightConfigFromFile("./light/testconfig.json"))

	os.Exit(0)

	cm := &CueMaster{}

	mainCueStack := CueStack{Priority: 2, Name: "main"}
	for x := 1; x <= 2; x++ {

		a := cm.NewCue([]CueFrame{
			cm.NewCueFrame([]CueFrameAction{
				cm.NewCueFrameAction(time.Millisecond*1500, RGBColor{}),
				cm.NewCueFrameAction(time.Second*time.Duration(x), RGBColor{}),
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
