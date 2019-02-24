package config

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

//Client represents client config
type Client struct {
	ServerAddress string
	Tracing       struct {
		Enabled       bool
		ServerAddress string
		ServiceName   string
	}
}

//GetClientConfig extracts Client config from context
func GetClientConfig(ctx context.Context) *Client {
	val := ctx.Value(ContextKeyClient)
	if val == nil {
		return nil
	}
	return val.(*Client)
}

//InjectIntoContext injects Client config into a provided ctx
func (s *Client) InjectIntoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyClient, s)
}

//LoadClient returns client config
func LoadClient() *Client {
	viper.SetConfigName("hyperion")
	viper.AddConfigPath("$HOME/.hyperion")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error config file: %s", err)
	}
	c := Client{}
	viper.SetDefault("client.server", "localhost:8888")
	c.ServerAddress = viper.GetString("client.server")

	if viper.IsSet("tracing") {
		c.Tracing.Enabled = true
		c.Tracing.ServerAddress = viper.GetString("tracing.server")
		c.Tracing.ServiceName = viper.GetString("tracing.servicename")
	}

	return &c
}
