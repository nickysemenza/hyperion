package light

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockOLAClient struct {
	mock.Mock
}

func (c *MockOLAClient) Close() {
	c.Called()
	return
}

func (c *MockOLAClient) SendDmx(universe int, values []byte) (status bool, err error) {
	args := c.Called(universe, values)
	return args.Bool(0), nil
}

func TestSendDMXWorker(t *testing.T) {
	client := new(MockOLAClient)
	client.On("Close")

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	m := Manager{
		dmxState: DMXState{universes: make(map[int][]byte)},
	}
	m.dmxState.set(ctx, dmxOperation{1, 1, 12})
	m.dmxState.set(ctx, dmxOperation{3, 9, 100})

	chans1 := make([]byte, 255)
	chans1[0] = 12
	client.On("SendDmx", 1, chans1).Return(true, nil)

	chans3 := make([]byte, 255)
	chans3[8] = 100
	client.On("SendDmx", 3, chans3).Return(true, nil)

	go SendDMXWorker(ctx, client, time.Millisecond*20, &m, &wg)
	time.Sleep(time.Millisecond * 100)
	cancel()
	wg.Wait()
	client.AssertExpectations(t)

}

func TestDMXAttributeChannels(t *testing.T) {
	tt := []struct {
		profile  config.LightProfileDMX
		name     string
		expected int
	}{
		{config.LightProfileDMX{Channels: map[string]int{"red": 1, "green": 2}}, "red", 1},
		{config.LightProfileDMX{Channels: map[string]int{"red": 1, "green": 2}}, "blue", 0},
	}
	for _, tc := range tt {
		res := getChannelIndexForAttribute(&tc.profile, tc.name)
		if res != tc.expected {
			t.Errorf("got channel index %d, expected %d", res, tc.expected)
		}
	}
}
func TestDMX(t *testing.T) {
	s1 := DMXState{universes: make(map[int][]byte)}
	s1.set(context.Background(), dmxOperation{2, 22, 40})

	require.EqualValues(t, 40, s1.universes[2][21], "didn't set DMX state instance properly")
	require.Error(t, s1.set(context.Background(), dmxOperation{2, 0, 2}), "should not allow channel 0")
}
