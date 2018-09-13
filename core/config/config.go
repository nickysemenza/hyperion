package config

import "context"

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
}

//GetClientConfig extracts Client config from context
func GetClientConfig(ctx context.Context) *Client {
	return ctx.Value(ContextKeyClient).(*Client)
}

//Client represents client config
type Client struct {
	ServerAddress string
}
