package cue

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/clock"
)

//MasterManager is an interface
type MasterManager interface {
	ProcessStack(ctx context.Context, cs *Stack, wg *sync.WaitGroup)
	ProcessCue(ctx context.Context, c *Cue, wg *sync.WaitGroup)
	ProcessFrame(ctx context.Context, cf *Frame, wg *sync.WaitGroup)
	ProcessFrameAction(ctx context.Context, cfa *FrameAction, wg *sync.WaitGroup)
	EnQueueCue(c Cue, cs *Stack) *Cue
	AddIDsRecursively(c *Cue)
	GetDefaultCueStack() *Stack
	ProcessForever(ctx context.Context, wg *sync.WaitGroup)
	GetLightManager() *light.Manager
}

//Master is the parent of all CueStacks, is a singleton
type Master struct {
	CueStacks    []Stack `json:"cue_stacks"`
	currentID    int64
	cl           clock.Clock
	idLock       sync.Mutex
	LightManager *light.Manager
}

//cueMaster singleton
var cueMasterSingleton Master

//InitializeMaster initializes the cuemaster
func InitializeMaster(cl clock.Clock, ls *light.Manager) MasterManager {
	return &Master{
		currentID:    1,
		cl:           cl,
		CueStacks:    []Stack{{Priority: 1, Name: "main"}},
		LightManager: ls,
	}
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

//ProcessForever runs all the cuestacks
func (cm *Master) ProcessForever(ctx context.Context, wg *sync.WaitGroup) {
	for x := range cm.CueStacks {
		wg.Add(1)
		go cm.ProcessStack(ctx, &cm.CueStacks[x], wg)
	}
}

//GetLightManager returns a poitner to the light state manager
func (cm *Master) GetLightManager() *light.Manager {
	return cm.LightManager
}
