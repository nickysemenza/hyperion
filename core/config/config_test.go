package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractConfigsFromContext(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	require.Nil(GetServerConfig(ctx))
	require.Nil(GetClientConfig(ctx))

	s := &Server{}
	s.Inputs.RPC.Enabled = true
	ctx = s.InjectIntoContext(ctx)
	require.EqualValues(s, GetServerConfig(ctx))

	c := &Client{}
	c.ServerAddress = "foo"
	ctx = c.InjectIntoContext(ctx)
	require.EqualValues(c, GetClientConfig(ctx))
}
