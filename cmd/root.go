package cmd

import (
	"fmt"
	"os"

	"github.com/nickysemenza/hyperion/cue"
	"github.com/nickysemenza/hyperion/light"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hyperion",
	Short: "Hyperion is hype",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Server")
		light.ReadLightConfigFromFile("./light/testconfig.json")

		go func() {
			c, _ := cue.BuildCueFromCommand("hue1:#00FF00:1000")
			cs := cue.GetCueMaster().GetDefaultCueStack()
			cs.EnQueueCue(*c)
		}()

		runServer()
	},
}

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Run the client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: client")
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
