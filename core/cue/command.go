package cue

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
)

var (
	errorWrongPartCount    = errors.New("command: wrong number of parts, should be lights:colors:timings")
	errorPartSizeMismatch  = errors.New("command: number of lights, colors, and timings must be the same")
	errorInvalidTime       = errors.New("command: invalid time")
	errorMissingFunction   = errors.New("command: fn(?) is missing")
	errorUndefinedFunction = errors.New("command: function is not defined")
	errorWrongPartLen      = errors.New("command: wrong number of colon delimited groups")
)

//NewFromCommand returns a cue based on a command.
func NewFromCommand(ctx context.Context, cmd string) (*Cue, error) {
	//remove spaces
	cmd = strings.Replace(cmd, " ", "", -1)

	//extracts: `match1(match2)`
	re := regexp.MustCompile(`(?m)(.*?)\((.*?)\)`)
	groups := re.FindAllStringSubmatch(cmd, -1)
	if len(groups) != 1 {
		return nil, errorMissingFunction
	}
	commandType := groups[0][1]
	argString := groups[0][2]

	args := strings.Split(argString, ",")

	var cue *Cue
	var err error
	switch commandType {

	case "set":
		cue, err = processSetCommand(argString)
	case "cycle":
		cue, err = processCycleCommand(argString)
	default:
		if userCommand, ok := config.GetServerConfig(ctx).Commands[commandType]; ok {
			cue, err = BuildCueFromUserCommand(ctx, userCommand, args)
		} else if systemCommand, ok := systemCommands[commandType]; ok {
			cue, err = BuildCueFromUserCommand(ctx, systemCommand, args)
		} else {
			err = errorUndefinedFunction
		}
	}

	if err != nil {
		return nil, err
	}
	return cue, nil

}

// e.g. cycle(c1+c2+c3+c4+c5+c6:500ms)
func processCycleCommand(cmd string) (*Cue, error) {
	cue := Cue{}
	parts := strings.Split(cmd, ":")
	if len(parts) != 2 {
		return nil, errorWrongPartLen
	}
	lightList := strings.Split(parts[0], "+")
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return nil, err
	}
	for x := range lightList {
		frame := Frame{}
		for y := 0; y < len(lightList); y++ {
			action := FrameAction{}
			action.LightName = lightList[y]

			action.NewState = light.TargetState{
				State:    light.State{RGB: color.GetRGBFromString("#0000FF")},
				Duration: duration,
			}
			if x == y {
				action.NewState = light.TargetState{
					State:    light.State{RGB: color.GetRGBFromString("#FF0000")},
					Duration: duration,
				}
			}

			frame.Actions = append(frame.Actions, action)
		}
		cue.Frames = append(cue.Frames, frame)
	}

	return &cue, nil
}

func processSetCommand(cmd string) (*Cue, error) {
	cue := Cue{}
	for _, cueFrameString := range strings.Split(cmd, "|") {
		frame := Frame{}
		parts := strings.Split(cueFrameString, ":")
		if len(parts) != 3 {
			return nil, errorWrongPartCount
		}
		if len(parts[0]) == 0 || len(parts[1]) == 0 || len(parts[2]) == 0 {
			return nil, errorWrongPartCount
		}
		lightList := strings.Split(parts[0], "+")
		colorList := strings.Split(parts[1], "+")
		timeList := strings.Split(parts[2], "+")

		numTimes := len(timeList)
		numColors := len(colorList)
		numLights := len(lightList)

		if !(numTimes == numColors && numColors == numLights) {
			return nil, errorPartSizeMismatch
		}
		for x := 0; x < numLights; x++ {
			action := FrameAction{}
			action.LightName = lightList[x]
			duration, err := time.ParseDuration(timeList[x])
			if err != nil {
				return nil, errorInvalidTime
			}
			action.NewState = light.TargetState{
				State:    light.State{RGB: color.GetRGBFromString(colorList[x])},
				Duration: duration,
			}
			frame.Actions = append(frame.Actions, action)
		}
		cue.Frames = append(cue.Frames, frame)
	}
	return &cue, nil
}
