package server

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/heatxsink/go-hue/lights"
	"github.com/nickysemenza/gola"
	"github.com/nickysemenza/hyperion/util/clock"

	log "github.com/sirupsen/logrus"

	"github.com/nickysemenza/hyperion/util/tracing"

	"github.com/nickysemenza/hyperion/api"
	"github.com/nickysemenza/hyperion/control/homekit"
	"github.com/nickysemenza/hyperion/control/job"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
)

//Run starts the server
func Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	c := config.GetServerConfig(ctx)

	//Initialize lights (including hue output)
	hueConn := lights.New(c.Outputs.Hue.Address, c.Outputs.Hue.Username)
	lm, err := light.NewManager(ctx, hueConn)
	if err != nil {
		log.Fatalf("error initializing light manager. err='%v'", err)
	}

	master := cue.InitializeMaster(clock.RealClock{}, lm)

	go tracing.InitTracer(ctx)

	//Set up Homekit Server
	wg.Add(1)
	go homekit.Start(ctx, &wg, master)

	//Set up RPC server
	wg.Add(1)
	go api.ServeRPC(ctx, &wg, master)

	//Setup API server
	wg.Add(1)
	go api.ServeHTTP(ctx, &wg, master)

	wg.Add(1)
	go job.ProcessForever(ctx, &wg, c.Jobs, master)

	//proceess cues forever
	master.ProcessForever(ctx, &wg)

	olaConfig := c.Outputs.OLA
	if !olaConfig.Enabled {
		log.Info("ola output is not enabled")
	} else {
		client, err := gola.New(olaConfig.Address)
		if err != nil {
			log.Errorf("could not start DMX worker: could not connect to ola: %v", err)
		} else {
			wg.Add(1)
			go light.SendDMXWorker(ctx, client, olaConfig.Tick, lm, &wg)
		}
	}

	//handle CTRL+C interrupt
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("shutting down hyperion server")
	cancel()
	wg.Wait()
}
