package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nickysemenza/hyperion/api"
	"github.com/nickysemenza/hyperion/control/homekit"
	"github.com/nickysemenza/hyperion/control/trigger"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/metrics"
)

func runServer(ctx context.Context) {
	metrics.Register()
	//Set up Homekit Server
	go homekit.Start(ctx)

	//Set up RPC server
	go api.ServeRPC(ctx)

	//Setup API server
	go api.ServeHTTP(ctx)

	//proceess cues forever
	cue.GetCueMaster().ProcessForever(ctx)

	go light.SendDMXWorker(ctx)

	//process triggers
	go trigger.ProcessTriggers(ctx)

	//handle CTRL+C
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Println("Shutdown hyperion ...")
		os.Exit(0)
	}()

	//keep going
	for {
		time.Sleep(1 * time.Second)
	}

}
