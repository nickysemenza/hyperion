package cue

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/clock"
	"github.com/nickysemenza/hyperion/util/color"
)

//MasterManager is an interface
type MasterManager interface {
	ProcessStack(ctx context.Context, cs *Stack)
	ProcessCue(ctx context.Context, c *Cue, wg *sync.WaitGroup)
	ProcessFrame(ctx context.Context, cf *Frame, wg *sync.WaitGroup)
	ProcessFrameAction(ctx context.Context, cfa *FrameAction, wg *sync.WaitGroup)
	EnQueueCue(c Cue, cs *Stack) *Cue
	AddIDsRecursively(c *Cue)
	GetDefaultCueStack() *Stack
	ProcessForever(ctx context.Context)
	GetLightManager() *light.StateManager
}

//Master is the parent of all CueStacks, is a singleton
type Master struct {
	CueStacks    []Stack `json:"cue_stacks"`
	currentID    int64
	cl           clock.Clock
	idLock       sync.Mutex
	LightManager *light.StateManager
}

//cueMaster singleton
var cueMasterSingleton Master

//InitializeMaster initializes the cuemaster
func InitializeMaster(cl clock.Clock, ls *light.StateManager) MasterManager {
	return &Master{
		currentID:    1,
		cl:           cl,
		CueStacks:    []Stack{{Priority: 1, Name: "main"}},
		LightManager: ls,
	}
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
	cm.idLock.Lock()
	defer cm.idLock.Unlock()

	id := cm.currentID
	cm.currentID++
	return id
}

//GetDefaultCueStack gives the first cuestack
func (cm *Master) GetDefaultCueStack() *Stack {
	return &cm.CueStacks[0]
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
		go cm.ProcessStack(ctx, &cm.CueStacks[x])
	}
}

//GetLightManager returns a poitner to the light state manager
func (cm *Master) GetLightManager() *light.StateManager {
	return cm.LightManager
}
