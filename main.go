package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nickysemenza/hyperion/client"
	"github.com/nickysemenza/hyperion/core/config"
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
		light.ReadLightConfigFromFile("lights.yaml") //TODO: move to viper setup

		server.Run(config.LoadServer().InjectIntoContext(context.Background()))
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
