package light

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
	s1 := getDMXStateInstance()
	s1.setDMXValues(context.Background(), dmxOperation{2, 22, 40})

	s2 := getDMXStateInstance()
	require.EqualValues(t, 40, s2.universes[2][21], "didn't set DMX state instance properly")
	require.Error(t, s2.setDMXValues(context.Background(), dmxOperation{2, 0, 2}), "should not allow channel 0")
	require.Equal(t, s1, s2, "should be a singleton!")
}

func TestDMXLight_blindlySetRGBToStateAndDMX(t *testing.T) {
	type fields struct {
		StartAddress int
		Universe     int
		Profile      string
	}
	tests := []struct {
		name   string
		fields fields
		color  color.RGB
	}{
		{"setLightToGreen", fields{Profile: "a", Universe: 4, StartAddress: 1}, color.GetRGBFromString("green")},
		{"withOffsetStartAddress", fields{Profile: "a", Universe: 4, StartAddress: 10}, color.GetRGBFromString("green")},
	}

	c := config.Server{}
	c.DMXProfiles = make(config.DMXProfileMap)
	c.DMXProfiles["a"] = config.LightProfileDMX{Name: "a", Channels: map[string]int{"red": 1, "green": 2, "blue": 3}}

	ctx := c.InjectIntoContext(context.Background())
	Initialize(ctx, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DMXLight{
				StartAddress: tt.fields.StartAddress,
				Universe:     tt.fields.Universe,
				Profile:      tt.fields.Profile,
			}
			d.blindlySetRGBToStateAndDMX(ctx, tt.color)
			ds := getDMXStateInstance()
			//green means first chan should be 0, secnd 255
			require.Equal(t, 0, ds.getDmxValue(tt.fields.Universe, tt.fields.StartAddress))
			require.Equal(t, 255, ds.getDmxValue(tt.fields.Universe, 1+tt.fields.StartAddress))
		})
	}
}

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
	go func() {
		time.Sleep(time.Second)
		cancel()

	}()

	s := getDMXStateInstance()
	s.universes = make(map[int][]byte) //TODO: make it so i don't have to reset
	s.setDMXValues(ctx, dmxOperation{1, 1, 12})
	s.setDMXValues(ctx, dmxOperation{3, 9, 100})

	chans1 := make([]byte, 255)
	chans1[0] = 12
	client.On("SendDmx", 1, chans1).Return(true, nil)

	chans3 := make([]byte, 255)
	chans3[8] = 100
	client.On("SendDmx", 3, chans3).Return(true, nil)

	go SendDMXWorker(ctx, client, time.Second, &wg)
	// spew.Dump(client.Calls)
	wg.Wait()
	client.AssertExpectations(t)

}
