package cue

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
)

//Master is the parent of all CueStacks, is a singleton
type Master struct {
	CueStacks []Stack `json:"cue_stacks"`
	currentID int64
	cl        Clock
	m         sync.Mutex
}

type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) Sleep(d time.Duration)                  { time.Sleep(d) }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

//cueMaster singleton
var (
	cueMasterSingleton Master
	once               sync.Once
)

//GetCueMaster makes a singleton for the cue master
func GetCueMaster() *Master {
	once.Do(func() {
		cueMasterSingleton = Master{currentID: 1, cl: realClock{}}
		cueMasterSingleton.CueStacks = append(cueMasterSingleton.CueStacks, cueMasterSingleton.NewStack(1, "main"))
	})
	return &cueMasterSingleton
}

//NewFrameAction creates a new instate with incr ID
func (cm *Master) NewFrameAction(duration time.Duration, color color.RGB, lightName string) FrameAction {
	return FrameAction{ID: cm.getNextIDForUse(), LightName: lightName, NewState: light.TargetState{State: light.State{RGB: color}, Duration: duration}}
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
	cm.m.Lock()
	defer cm.m.Unlock()

	id := cm.currentID
	cm.currentID++
	return id
}

//GetDefaultCueStack gives the first cuestack
func (cm *Master) GetDefaultCueStack() *Stack {
	return &cm.CueStacks[0]
}

//NewStack makes a new cue stack
func (cm *Master) NewStack(priority int, name string) Stack {
	return Stack{Priority: priority, Name: name}
}

//NewFrame creates a new instate with incr ID
func (cm *Master) NewFrame(actions []FrameAction) Frame {
	return Frame{ID: cm.getNextIDForUse(), Actions: actions}
}

//New creates a new instate with incr ID
func (cm *Master) New(frames []Frame, name string) Cue {
	return Cue{ID: cm.getNextIDForUse(), Frames: frames, Status: statusEnqueued}
}

//ProcessForever runs all the cuestacks
func (cm *Master) ProcessForever(ctx context.Context) {
	for x := range cm.CueStacks {
		go cm.CueStacks[x].ProcessStack(ctx)
	}
}
