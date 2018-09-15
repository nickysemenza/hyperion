package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nickysemenza/hyperion/client"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		viper.SetConfigName("config")          // name of config file (without extension)
		viper.AddConfigPath("$HOME/.hyperion") // call multiple times to add many search paths
		viper.AddConfigPath(".")               // optionally look for config in the working directory
		err := viper.ReadInConfig()            // Find and read the config file
		if err != nil {                        // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s \n", err))
		}

		viper.Debug()

		c := config.Server{}
		//inputs
		if viper.IsSet("inputs.rpc") {
			c.Inputs.RPCAddress = viper.GetString("inputs.rpc.address")
		}
		if viper.IsSet("inputs.http") {
			c.Inputs.HTTPAddress = viper.GetString("inputs.http.address")
		}

		//outputs
		if viper.IsSet("outputs.ola") {
			c.Outputs.OLA.Enabled = true
			c.Outputs.OLA.Address = viper.GetString("outputs.ola.address")
		}
		if viper.IsSet("outputs.hue") {
			c.Outputs.Hue.Enabled = true
			c.Outputs.Hue.Address = viper.GetString("outputs.hue.address")
			c.Outputs.Hue.Username = viper.GetString("outputs.hue.username")
		}

		//other
		if viper.IsSet("outputs.hue") {
			c.Tracing.Enabled = true
			c.Tracing.ServerAddress = viper.GetString("outputs.tracing.server")
			c.Tracing.ServiceName = viper.GetString("outputs.tracing.servicename")
		}

		server.Run(context.WithValue(context.Background(), config.ContextKeyServer, &c))
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
