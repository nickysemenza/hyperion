package server

import (
	"context"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
)

func TestServerBasic(t *testing.T) {
	config := config.Server{}
	ctx, cancel := context.WithCancel(config.InjectIntoContext(context.Background()))
	go Run(ctx)
	time.Sleep(time.Millisecond * 100)
	cancel()
}
