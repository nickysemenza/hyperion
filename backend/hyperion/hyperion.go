package hyperion

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nickysemenza/hyperion/backend/api"
	"github.com/nickysemenza/hyperion/backend/cue"
	"github.com/nickysemenza/hyperion/backend/homekit"
	"github.com/nickysemenza/hyperion/backend/light"
	"github.com/nickysemenza/hyperion/backend/trigger"
)

func RunServer() {
	//Set up Homekit Server
	go homekit.Start()

	//Set up RPC server
	go api.ServeRPC(8888)

	//Setup API server
	go api.ServeHTTP()

	//proceess cues forever
	cue.GetCueMaster().ProcessForever()

	go light.SendDMXValuesToOLA()

	//process triggers
	go trigger.ProcessTriggers()

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
