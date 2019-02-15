//Package job is used for cron scheduling of command execution
package job

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/robfig/cron"
)

//ProcessForever begins the cron job runner
func ProcessForever(ctx context.Context, wg *sync.WaitGroup, jobs []config.Job, master cue.MasterManager) {
	log.Println("starting job worker")
	defer wg.Done()
	c := cron.New()
	for _, x := range jobs {
		job := x
		c.AddFunc(job.Cron, func() {
			c, err := cue.CommandToCue(ctx, master, job.Command)
			if err != nil {
				log.Printf("failed to build command for job, job='%v', error='%v'", job, err)
			} else {
				stack := master.GetDefaultCueStack()
				c.Source.Input = cue.SourceInputJob
				c.Source.Type = cue.SourceTypeCommand
				c.Source.Meta = fmt.Sprintf("job=%v", job)
				master.EnQueueCue(ctx, *c, stack)
			}
			fmt.Printf("processing cron: %s", job.Command)
		})

	}
	c.Start()

	<-ctx.Done()
	log.Println("job worker stopped")

}
