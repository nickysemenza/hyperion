package config

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

//Client represents client config
type Client struct {
	ServerAddress string
}

//GetClientConfig extracts Client config from context
func GetClientConfig(ctx context.Context) *Client {
	return ctx.Value(ContextKeyClient).(*Client)
}

//InjectIntoContext injects Client config into a provided ctx
func (s *Client) InjectIntoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyClient, s)
}

//LoadClient returns client config
func LoadClient() *Client {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.hyperion")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error config file: %s", err)
	}
	c := Client{}
	viper.SetDefault("client.server", "localhost:8888")
	c.ServerAddress = viper.GetString("client.server")
	return &c
}
