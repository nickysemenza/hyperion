package light

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/require"

	"github.com/heatxsink/go-hue/hue"
	"github.com/heatxsink/go-hue/lights"
	"github.com/stretchr/testify/mock"
)

func TestGetTransitionTimeAs100msMultiple(t *testing.T) {
	tt := []struct {
		input    time.Duration
		expected uint16
	}{
		{time.Duration(time.Second), 10},
		{time.Duration(time.Millisecond * 200), 2},
		{time.Duration(time.Millisecond * 250), 2},
		{time.Duration(time.Millisecond * 270), 2},
		{time.Duration(0), 0},
	}
	for _, x := range tt {
		res := getTransitionTimeAs100msMultiple(x.input)
		if res != x.expected {
			t.Errorf("got %d, expected %d", res, x.expected)
		}
	}
}

type MockHueConn struct {
	mock.Mock
}

func (h *MockHueConn) SetLightState(lightID int, state lights.State) ([]hue.ApiResponse, error) {
	h.Called(lightID, state)
	return nil, nil
}

func (h *MockHueConn) GetAllLights() ([]lights.Light, error) {
	args := h.Called()
	return args.Get(0).([]lights.Light), args.Error(1)
}

func TestGetDiscoveredHues(t *testing.T) {
	tests := []struct {
		name     string
		resp     []lights.Light
		err      error
		expected map[string]int
	}{
		{
			name:     "success",
			resp:     []lights.Light{{ID: 4, Name: "foo"}},
			err:      nil,
			expected: map[string]int{"foo": 4},
		},
		{
			name:     "fail",
			resp:     []lights.Light{},
			err:      fmt.Errorf("foo"),
			expected: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := new(MockHueConn)
			m := StateManager{
				hueConnection: h,
			}
			h.On("GetAllLights").Return(tt.resp, tt.err)
			require.Equal(t, tt.expected, m.GetDiscoveredHues().ByName)
			h.AssertExpectations(t)
		})
	}
}

func TestSetColor(t *testing.T) {

	tests := []struct {
		name                   string
		color                  color.RGB
		timing                 time.Duration
		expectedBrightness     uint8
		expectedIsOn           bool
		expectedTransitionTime uint16
	}{
		{
			name:                   "set to blue",
			color:                  color.RGB{B: 255},
			timing:                 time.Second,
			expectedBrightness:     255,
			expectedIsOn:           true,
			expectedTransitionTime: 10,
		},
		{
			name:                   "set to black",
			color:                  color.RGB{},
			timing:                 time.Second / 10,
			expectedBrightness:     0,
			expectedIsOn:           false,
			expectedTransitionTime: 1,
		},
		{
			name:                   "set to white",
			color:                  color.RGB{R: 255, G: 255, B: 255},
			timing:                 time.Millisecond * 250,
			expectedBrightness:     255,
			expectedIsOn:           true,
			expectedTransitionTime: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := HueLight{
				HueID: 4,
				Name:  "foo",
			}
			s := &config.Server{}
			ctx := s.InjectIntoContext(context.Background())

			h := new(MockHueConn)
			sm, err := NewManager(ctx, h)
			require.NoError(t, err)
			h.On("SetLightState", hl.HueID, mock.AnythingOfType("lights.State"))

			targetState := TargetState{Duration: tt.timing}
			targetState.State.RGB = tt.color
			hl.SetState(context.Background(), sm, targetState)
			time.Sleep(time.Millisecond)
			h.AssertExpectations(t)

			require.Len(t, h.Calls, 1)
			stateParam := h.Calls[0].Arguments[1].(lights.State)
			require.Equal(t, tt.expectedIsOn, stateParam.On)
			require.Equal(t, tt.expectedBrightness, stateParam.Bri)
			require.Equal(t, tt.expectedTransitionTime, stateParam.TransitionTime)
			// spew.Dump(h.Calls)

		})
	}
}
