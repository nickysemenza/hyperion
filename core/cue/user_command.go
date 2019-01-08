package cue

import (
	"context"
	"fmt"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

//LuaAction is a minimal version of Action that lua code returns
type LuaAction struct {
	Light  string
	Color  string
	Timing string
}

//LuaFrame is a minimal version of Frame that lua code returns
type LuaFrame struct {
	Actions []LuaAction
}

//LuaCue is a minimal version of Cue that lua code returns
type LuaCue struct {
	Frames []LuaFrame
}

//toCue builds a real Cue from the LuaCue
func (lc *LuaCue) toCue() (*Cue, error) {
	cue := Cue{}
	for _, f := range lc.Frames {
		frame := Frame{}
		for _, fa := range f.Actions {
			action := FrameAction{}
			duration, err := time.ParseDuration(fa.Timing)
			if err != nil {
				return nil, err
			}
			action.LightName = fa.Light
			action.NewState = light.TargetState{
				State:    light.State{RGB: color.GetRGBFromString(fa.Color)},
				Duration: duration,
			}
			frame.Actions = append(frame.Actions, action)
		}
		cue.Frames = append(cue.Frames, frame)
	}
	return &cue, nil
}

//LuaToHex is a lua helper to convert rgb to hex
func LuaToHex(L *lua.LState) int {
	r := L.ToInt(1)
	g := L.ToInt(2)
	b := L.ToInt(3)
	rgb := color.RGB{R: r, G: g, B: b}
	L.Push(lua.LString(rgb.ToHex()))
	return 1
}

//BuildCueFromUserCommand processes a lua user command
func BuildCueFromUserCommand(ctx context.Context, m MasterManager, command config.UserCommand, args []string) (*Cue, error) {
	L := lua.NewState()
	defer L.Close()

	//build lua list of light names
	luaLightList := L.NewTable()

	for _, name := range m.GetLightManager().GetLightNames() {
		luaLightList.Append(lua.LString(name))
	}

	L.SetGlobal("light_list", luaLightList)
	L.SetGlobal("rgb_to_hex", L.NewFunction(LuaToHex))

	//transform arg list
	lArgs := make([]lua.LValue, len(args))
	for x := range args {
		lArgs[x] = lua.LString(args[x])
	}

	//run the lua command blob
	if err := L.DoString(command.Body); err != nil {
		return nil, fmt.Errorf("user command processing error, could not run provided lua code err=%v", err)
	}
	//call user definedprocess func
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("process"),
		NRet:    1,
		Protect: true,
	}, lArgs...); err != nil {
		return nil, fmt.Errorf("user command processing error, could not call process() err=%v", err)
	}
	//get/pop the return value
	ret := L.Get(-1)
	L.Pop(1)

	//unmarshal returned data into a LuaCue
	var c LuaCue
	if err := gluamapper.Map(ret.(*lua.LTable), &c); err != nil {
		return nil, fmt.Errorf("user command processing error could not unmarshal result, err=%v", err)
	}

	return c.toCue()
}
