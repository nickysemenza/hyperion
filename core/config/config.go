package config

import (
	"context"
	"fmt"
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
	RPCAddress  string
	HTTPAddress string
	Outputs     struct {
		OLA struct {
			Address string
		}
		Hue struct {
			Address  string
			Username string
		}
	}
	Tracing struct {
		ServerAddress string
		ServiceName   string
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
