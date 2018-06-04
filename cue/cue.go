package cue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nickysemenza/hyperion/color"
	"github.com/nickysemenza/hyperion/light"
	log "github.com/sirupsen/logrus"
)

type ctxKey int

const (
	keyStackName ctxKey = iota
	keyCueID
	keyFrameID
	keyFrameActionID
)

func getLogrusFieldsFromContext(ctx context.Context) log.Fields {
	return log.Fields{
		"action_id":  ctx.Value(keyFrameActionID),
		"frame_id":   ctx.Value(keyFrameID),
		"cue_id":     ctx.Value(keyCueID),
		"stack_name": ctx.Value(keyStackName),
	}
}

//Stack is basically a precedence priority queue (really a CueQueue sigh)
type Stack struct {
	Priority      int    `json:"priority"`
	Name          string `json:"name"`
	Cues          []Cue  `json:"cues"`
	ProcessedCues []Cue  `json:"processed_cues"`
	m             sync.Mutex
	ActiveCue     *Cue `json:"active_cue"`
}

//Cue is a cue.
type Cue struct {
	ID              int64   `json:"id"`
	Frames          []Frame `json:"frames"`
	Name            string  `json:"name"`
	shouldRepeat    bool
	shouldHoldAfter bool          //default false, will pause the CueStack after executing this cue, won't move on to next
	WaitBefore      time.Duration `json:"wait_before"`
	WaitAfter       time.Duration `json:"wait_after"`
	StartedAt       time.Time     `json:"started_at"`
	FinishedAt      time.Time     `json:"finished_at"`
	RealDuration    time.Duration `json:"real_duration"`
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

//ProcessStack processes cues
func (cs *Stack) ProcessStack() {
	log.Printf("[CueStack: %s]\n", cs.Name)
	for {
		if nextCue := cs.deQueueNextCue(); nextCue != nil {
			ctx := context.WithValue(context.Background(), keyStackName, cs.Name)
			cs.ActiveCue = nextCue
			nextCue.StartedAt = time.Now()
			nextCue.ProcessCue(ctx)
			nextCue.FinishedAt = time.Now()
			nextCue.RealDuration = nextCue.FinishedAt.Sub(nextCue.StartedAt)
			cs.ActiveCue = nil
			cs.ProcessedCues = append(cs.ProcessedCues, *nextCue)
		} else {
			// fmt.Println("FINISHED PROCESSING CUESTACK, RESTARTING")
		}
	}
}

func (cs *Stack) deQueueNextCue() *Cue {
	if len(cs.Cues) > 0 {
		cs.m.Lock()
		x := cs.Cues[0]
		cs.Cues = cs.Cues[1:]
		cs.m.Unlock()
		return &x
	}
	return nil
}

//EnQueueCue puts a cue on the queue
func (cs *Stack) EnQueueCue(c Cue) {
	cs.m.Lock()
	defer cs.m.Unlock()
	log.WithFields(log.Fields{"cue_id": c.ID, "stack_name": cs.Name}).Info("enqueued!")

	cs.Cues = append(cs.Cues, c)
}

//ProcessCue processes cue
func (c *Cue) ProcessCue(ctx context.Context) {
	ctx = context.WithValue(ctx, keyCueID, c.ID)
	log.WithFields(getLogrusFieldsFromContext(ctx)).Info("ProcessCue")
	for _, eachFrame := range c.Frames {
		eachFrame.ProcessFrame(ctx)
	}
}

//AddIDsRecursively populates the ID fields on a cue, its frames, and their actions
func (c *Cue) AddIDsRecursively() {
	cm := GetCueMaster()
	if c.ID == 0 {
		c.ID = cm.getNextIDForUse()
	}
	for x := range c.Frames {
		eachFrame := &c.Frames[x]
		if eachFrame.ID == 0 {
			eachFrame.ID = cm.getNextIDForUse()
		}
		for y := range eachFrame.Actions {
			eachAction := &eachFrame.Actions[y]
			if eachAction.ID == 0 {
				eachAction.ID = cm.getNextIDForUse()
			}
		}
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
func (cf *Frame) ProcessFrame(ctx context.Context) {
	ctx = context.WithValue(ctx, keyFrameID, cf.ID)
	log.WithFields(getLogrusFieldsFromContext(ctx)).
		WithFields(log.Fields{"duration": cf.GetDuration(), "num_actions": len(cf.Actions)}).
		Info("ProcessFrame")

	// fmt.Println(cf.Actions)
	for x := range cf.Actions {
		go cf.Actions[x].ProcessFrameAction(ctx)
	}
	//no blocking, so wait until all the child frames have theoretically finished
	time.Sleep(cf.GetDuration())
}

//ProcessFrameAction does job stuff
func (cfa *FrameAction) ProcessFrameAction(ctx context.Context) {
	ctx = context.WithValue(ctx, keyFrameActionID, cfa.ID)
	now := time.Now().UnixNano() / int64(time.Millisecond)
	log.WithFields(getLogrusFieldsFromContext(ctx)).
		WithFields(log.Fields{"duration": cfa.NewState.Duration, "now_ms": now, "light": cfa.LightName}).
		Infof("ProcessFrameAction (color=%v)", cfa.NewState.RGB.String())

	if l := light.GetByName(cfa.LightName); l != nil {
		go l.SetState(ctx, cfa.NewState)
	} else {
		fmt.Printf("Cannot find light by name: %s\n", cfa.LightName)
	}
	//goroutine doesn't block, so hold until the SetState has (hopefully) finished timing-wise
	time.Sleep(cfa.NewState.Duration)
}

//MarshalJSON override that injects the full Light object.
func (cfa *FrameAction) MarshalJSON() ([]byte, error) {
	type Alias FrameAction
	return json.Marshal(&struct {
		Light light.Light `json:"light"`
		*Alias
	}{
		Light: light.GetByName(cfa.LightName),
		Alias: (*Alias)(cfa),
	})
}

//NewSimple returns a Cue that transitions the given light to the given color
func NewSimple(lightName string, c color.RGB) Cue {
	cm := GetCueMaster()
	return cm.New([]Frame{
		cm.NewFrame([]FrameAction{
			cm.NewFrameAction(time.Millisecond*500, c, lightName),
		})}, "")

}
