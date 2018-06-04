package cue

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	"github.com/nickysemenza/hyperion/backend/color"
	"github.com/nickysemenza/hyperion/backend/light"
)

//Master is the parent of all CueStacks, is a singleton
type Master struct {
	CueStacks []Stack `json:"cue_stacks"`
	currentID int64
	m         sync.Mutex
}

//cueMaster singleton
var (
	cueMasterSingleton Master
	once               sync.Once
)

//GetCueMaster makes a singleton for the cue master
func GetCueMaster() *Master {
	once.Do(func() {
		cueMasterSingleton = Master{currentID: 1}
	})
	return &cueMasterSingleton
}

//NewFrameAction creates a new instate with incr ID
func (cm *Master) NewFrameAction(duration time.Duration, color color.RGB, lightName string) FrameAction {
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
	cm.m.Lock()
	defer cm.m.Unlock()

	id := cm.currentID
	cm.currentID++
	return id
}

//GetDefaultCueStack gives the first cuestack
func (cm *Master) GetDefaultCueStack() *Stack {
	cm.m.Lock()
	defer cm.m.Unlock()
	if len(cm.CueStacks) == 0 {
		cm.CueStacks = append(cm.CueStacks, cm.NewStack(2, "main"))
	}
	return &cm.CueStacks[0]
}

//NewStack makes a new cue stack
func (cm *Master) NewStack(priority int, name string) Stack {
	return Stack{Priority: 2, Name: "main"}
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
