package config

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type ctxKey int

//Context keys
const (
	ContextKeyServer ctxKey = iota
	ContextKeyClient
)

//GetServerConfig extracts Server config from context
func GetServerConfig(ctx context.Context) *Server {
	return ctx.Value(ContextKeyServer).(*Server)
}

//InjectIntoContext injects Server config into a provided ctx
func (s *Server) InjectIntoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyServer, s)
}

//Server represents server config
type Server struct {
	Inputs struct {
		RPC struct {
			Enabled bool
			Address string
		}
		HTTP struct {
			Enabled        bool
			Address        string
			WSTickInterval time.Duration
		}
		HomeKit struct {
			Enabled bool
			Pin     string
		}
	}
	Outputs struct {
		OLA struct {
			Enabled bool
			Address string
			Tick    time.Duration
		}
		Hue struct {
			Enabled  bool
			Address  string
			Username string
		}
	}
	Tracing struct {
		Enabled       bool
		ServerAddress string
		ServiceName   string
	}

	Timings struct {
		FadeInterpolationTick time.Duration
		CueBackoff            time.Duration
	}
	Triggers []Trigger

	Lights struct {
		Hue     []LightHue
		DMX     []LightDMX
		Generic []LightGeneric
	}
	DMXProfiles DMXProfileMap
	Commands    UserCommandMap
}

//DMXProfileMap represents a map of dmx profiles
type DMXProfileMap map[string]LightProfileDMX

//LightHue holds config info for a Hue
type LightHue struct {
	Name  string
	HueID int `mapstructure:"hue_id"`
}

//LightGeneric holds config info for a Generic
type LightGeneric struct {
	Name string
}

//LightDMX hol;ds config info for a dmx light
type LightDMX struct {
	Name         string
	StartAddress int `mapstructure:"start_address"`
	Universe     int
	Profile      string
}

//LightProfileDMX holds config info for a dmx profile: channel and capability mappings
type LightProfileDMX struct {
	Name         string
	Capabilities []string
	Channels     map[string]int
}

//UserCommand holds a user command
type UserCommand struct {
	Body string
}

//UserCommandMap is used to map user commands by name
type UserCommandMap map[string]UserCommand

//Trigger holds configuration for a trigger
type Trigger struct {
	ID      int
	Source  string
	Command string
}

var (
	//GitCommit is injected at build time with the commit hash
	GitCommit = "0"
)

//GetVersion returns the currently running version
func GetVersion() string {
	return fmt.Sprintf("git-%s", GitCommit)
}

//LoadServer returns the server config (using viper)
func LoadServer() *Server {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.hyperion")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error config file: %s", err)
	}

	c := Server{}
	//inputs
	if viper.IsSet("inputs.rpc") {
		c.Inputs.RPC.Enabled = true
		c.Inputs.RPC.Address = viper.GetString("inputs.rpc.address")
	}
	if viper.IsSet("inputs.http") {
		c.Inputs.HTTP.Enabled = true
		c.Inputs.HTTP.Address = viper.GetString("inputs.http.address")
		viper.SetDefault("inputs.http.ws-tick", time.Millisecond*50)
		c.Inputs.HTTP.WSTickInterval = viper.GetDuration("inputs.http.ws-tick")
	}
	if viper.IsSet("inputs.homekit") {
		c.Inputs.HomeKit.Enabled = true
		viper.SetDefault("inputs.homekit.pin", "10000000")
		c.Inputs.HomeKit.Pin = viper.GetString("inputs.homekit.pin")
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

	//triggers
	err = viper.UnmarshalKey("triggers", &c.Triggers)

	//light config
	viper.UnmarshalKey("lights", &c.Lights)
	viper.UnmarshalKey("dmx_profiles", &c.DMXProfiles)
	viper.UnmarshalKey("commands", &c.Commands)
	spew.Dump(c)

	return &c
}
