package config

import (
	"context"
	"fmt"
	"time"
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
}

//GetClientConfig extracts Client config from context
func GetClientConfig(ctx context.Context) *Client {
	return ctx.Value(ContextKeyClient).(*Client)
}

//Client represents client config
type Client struct {
	ServerAddress string
}

var (
	//GitCommit is injected at build time with the commit hash
	GitCommit = "0"
)

func GetVersion() string {
	return fmt.Sprintf("git-%s", GitCommit)
}
