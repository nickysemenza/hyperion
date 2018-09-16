package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.hyperion")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}

		viper.Debug()

		c := config.Server{}
		//inputs
		if viper.IsSet("inputs.rpc") {
			c.Inputs.RPC.Enabled = true
			c.Inputs.RPC.Address = viper.GetString("inputs.rpc.address")
		}
		if viper.IsSet("inputs.http") {
			c.Inputs.HTTP.Enabled = true
			c.Inputs.HTTP.Address = viper.GetString("inputs.http.address")
			viper.SetDefault("inputs.http.ws-tick", time.Millisecond*50)
			c.Outputs.OLA.Tick = viper.GetDuration("inputs.http.ws-tick")
		}

		//outputs
		if viper.IsSet("outputs.ola") {
			c.Outputs.OLA.Enabled = true
			c.Outputs.OLA.Address = viper.GetString("outputs.ola.address")
			viper.SetDefault("outputs.ola.tick", time.Millisecond*50)
			c.Outputs.OLA.Tick = viper.GetDuration("outputs.ola.tick")
		}
		if viper.IsSet("outputs.hue") {
			c.Outputs.Hue.Enabled = true
			c.Outputs.Hue.Address = viper.GetString("outputs.hue.address")
			c.Outputs.Hue.Username = viper.GetString("outputs.hue.username")
		}

		//other
		if viper.IsSet("tracing") {
			c.Tracing.Enabled = true
			c.Tracing.ServerAddress = viper.GetString("tracing.server")
			c.Tracing.ServiceName = viper.GetString("tracing.servicename")
		}

		//timings
		viper.SetDefault("timings.fade-interpolation-tick", time.Millisecond*25)
		c.Timings.FadeInterpolationTick = viper.GetDuration("timings.fade-interpolation-tick")
		viper.SetDefault("timings.fade-cue-backoff", time.Millisecond*25)
		c.Timings.CueBackoff = viper.GetDuration("timings.fade-cue-backoff")

		// spew.Dump(c)
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
