package trigger

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"

	log "github.com/sirupsen/logrus"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
)

type trigger struct {
	source string
	id     int
}

//Action is called when an trigger needs to be fired
func Action(ctx context.Context, source string, id int) {
	process(ctx, trigger{source, id})
}

func process(ctx context.Context, t trigger) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "process trigger")
	defer span.Finish()
	span.LogKV("trigger", t)
	triggerConf := config.GetServerConfig(ctx).Triggers
	log.Printf("new trigger! %v\n", t)
	for _, each := range triggerConf {
		if each.ID == t.id && each.Source == t.source {
			if c, err := cue.NewFromCommand(ctx, each.Command); err != nil {
				log.Errorf("failed to build command from trigger, trigger=%v, command=%v", t, each.Command)
			} else {
				stack := cue.GetCueMaster().GetDefaultCueStack()
				c.Source.Input = cue.SourceInputTrigger
				c.Source.Type = cue.SourceTypeCommand
				c.Source.Meta = fmt.Sprintf("trigger=%s:%d", t.source, t.id)
				stack.EnQueueCue(*c)
			}
			// TODO: require one command per trigger, return here
		}
	}
}
