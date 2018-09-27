package cue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/stretchr/testify/require"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaToHex(t *testing.T) {

	tests := []struct {
		r    int
		g    int
		b    int
		want string
	}{
		{255, 120, 0, "#ff7800"},
		{0, 0, 0, "#000000"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		testName := fmt.Sprintf("rgb: (%d,%d,%d) -> %s", tt.r, tt.g, tt.b, tt.want)
		t.Run(testName, func(t *testing.T) {
			ls := lua.NewState()
			ls.Push(lua.LNumber(tt.r)) //r
			ls.Push(lua.LNumber(tt.g)) //g
			ls.Push(lua.LNumber(tt.b)) //b

			LuaToHex(ls)

			res := ls.Get(4).String()
			require.Equal(t, tt.want, res)
		})
	}
}

func TestToCue(t *testing.T) {
	tests := []struct {
		lc      *LuaCue
		c       *Cue
		wantErr bool
	}{
		{&LuaCue{}, &Cue{}, false},
		{&LuaCue{
			Frames: []LuaFrame{
				{Actions: []LuaAction{
					{Light: "foo", Timing: "2s", Color: "blue"},
				}},
			},
		}, &Cue{
			Frames: []Frame{
				{Actions: []FrameAction{
					{LightName: "foo",
						NewState: light.TargetState{
							Duration: time.Second * 2,
							State:    light.State{RGB: color.RGB{B: 255}},
						}},
				}},
			},
		}, false},
		{&LuaCue{
			Frames: []LuaFrame{
				{Actions: []LuaAction{
					{Light: "foo", Timing: "aa", Color: "blue"},
				}},
			},
		}, nil, true},
	}
	for _, tt := range tests {
		cue, err := tt.lc.toCue()
		if tt.wantErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		require.EqualValues(t, tt.c, cue)

	}
}

func TestBuildCueFromUserCommand(t *testing.T) {

	light.SetCurrentState("light2", light.State{})
	tests := []struct {
		name    string
		command config.UserCommand
		want    *Cue
		wantErr bool
	}{
		{"nil cmd",
			config.UserCommand{}, nil, true},
		{"bad user code",
			config.UserCommand{Body: "badcode("},
			nil, true},
		{"bad user code return val", config.UserCommand{Body: `
		function process()
			return {"bad structed return val"}
		end
				`}, nil, true},
		{"basic command",
			config.UserCommand{Body: `
			function process()
				action1 = {
					light = "light1",
					color = "blue",
					timing = "3s"
				}
				frame1 = {
					actions = {action1}
				}
				return {
					frames = {frame1}
				}
			end
		`}, &Cue{
				Frames: []Frame{
					{Actions: []FrameAction{
						{LightName: "light1",
							NewState: light.TargetState{
								Duration: time.Second * 3,
								State:    light.State{RGB: color.RGB{B: 255}},
							}},
					}},
				},
			}, false},
		{"command using light_list",
			config.UserCommand{Body: `
			function process()
				actions = {}
				for i, light in ipairs(light_list) do
					table.insert(actions, {
					light = light,
					color = rgb_to_hex(100,0,120),
					timing = "1s"
					})
				end
				frame1 = {
					actions = actions
				}
				return {
					frames = {frame1}
				}
			end
		`}, &Cue{
				Frames: []Frame{
					{Actions: []FrameAction{
						{LightName: "light2",
							NewState: light.TargetState{
								Duration: time.Second,
								State:    light.State{RGB: color.RGB{R: 100, B: 120}},
							}},
					}},
				},
			}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cue, err := BuildCueFromUserCommand(ctx, tt.command, "")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tt.want, cue)
		})
	}
}
