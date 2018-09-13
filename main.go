package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/hyperion/client"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/server"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdClient)
	rootCmd.AddCommand(cmdVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "hyperion",
	Short:   "Hyperion lighting controller v0.1",
	Version: config.GetVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}
var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Server, version:" + config.GetVersion())
		light.ReadLightConfigFromFile("./core/light/testconfig.yaml")

		go func() {
			c, _ := cue.BuildCueFromCommand("hue1:#00FF00:1000")
			cs := cue.GetCueMaster().GetDefaultCueStack()
			cs.EnQueueCue(*c)
		}()

		//TODO: read from file
		c := config.Server{
			RPCAddress:  ":8888",
			HTTPAddress: ":8080",
		}
		c.Outputs.OLA.Address = "localhost:9010"
		c.Outputs.Hue.Address = "10.0.0.39"
		c.Outputs.Hue.Username = "alW0LsA1mnXB28T4txGs01BeHi1WBr661VZ1eqEF"

		ctx := context.WithValue(context.Background(), config.ContextKeyServer, &c)
		spew.Dump(ctx)
		server.Run(ctx)
	},
}

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Run the client",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.Client{
			ServerAddress: "localhost:8888",
		}

		ctx := context.WithValue(context.Background(), config.ContextKeyClient, &c)
		client.Run(ctx)
	},
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Get version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.GetVersion())
	},
}
