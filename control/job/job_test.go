package job

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
)

func TestProcessForever(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m := new(cue.MockMasterManager)
	cs := cue.Stack{Name: "foo"}
	m.On("GetDefaultCueStack").Return(&cs)
	newCue, err := cue.NewFromCommand(ctx, m, "set(light1:#FF00FF:0)")
	newCue.Source = cue.Source{
		Input: cue.SourceInputJob,
		Type:  cue.SourceTypeCommand,
		Meta:  "job={foo set(light1:#FF00FF:0) * * * * *}",
	}
	require.NoError(t, err)
	m.On("EnQueueCue", *newCue, &cs).Return(newCue)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go ProcessForever(ctx, &wg, []config.Job{
		{Name: "foo", Command: "set(light1:#FF00FF:0)", Cron: "* * * * *"},
		{Name: "foo2bad", Command: "asdf", Cron: "* * * * *"},
	}, m)
	time.Sleep(time.Second * 3)
	cancel()
	wg.Wait()
	assert.Equal(t, 1, 1)
	m.AssertExpectations(t)
}
