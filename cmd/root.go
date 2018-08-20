package cmd

import (
	"fmt"
	"os"

	"github.com/nickysemenza/hyperion/client"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/spf13/cobra"
)

type cliHandler interface {
	startServer()
	startClient()
}

type defaultHandler struct{}

func (d defaultHandler) startServer() {
	runServer()
}
func (d defaultHandler) startClient() {
	client.Run("localhost:8888")
}

type cliConfig struct {
	h cliHandler
}

var mainHandler = &cliConfig{h: &defaultHandler{}}

var rootCmd = &cobra.Command{
	Use:   "hyperion",
	Short: "Hyperion lighting controller v0.1",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}
var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Server")
		light.ReadLightConfigFromFile("./core/light/testconfig.yaml")

		go func() {
			c, _ := cue.BuildCueFromCommand("hue1:#00FF00:1000")
			cs := cue.GetCueMaster().GetDefaultCueStack()
			cs.EnQueueCue(*c)
		}()

		// runServer()
		mainHandler.h.startServer()
	},
}

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Run the client",
	Run: func(cmd *cobra.Command, args []string) {
		mainHandler.h.startClient()
	},
}

//Execute gives control to cobra
func Execute() {

	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdClient)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
