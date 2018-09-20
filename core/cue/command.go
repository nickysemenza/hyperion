package cue

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
)

const (
	commandErrorWrongPartCount    = "command: wrong number of parts, should be lights:colors:timings"
	commandErrorPartSizeMismatch  = "command: number of lights, colors, and timings must be the same"
	commandErrorInvalidTime       = "command: invalid time"
	commandErrorMissingFunction   = "command: fn(?) is missing"
	commandErrorUndefinedFunction = "command: function is not defined"
)

//NewFromCommand returns a cue based on a command.
func NewFromCommand(cmd string) (*Cue, error) {
	cmd = strings.Replace(cmd, " ", "", -1)

	//extracts: `match1(match2)`
	re := regexp.MustCompile(`(?m)(.*?)\((.*?)\)`)
	groups := re.FindAllStringSubmatch(cmd, -1)
	if len(groups) != 1 {
		return nil, errors.New(commandErrorMissingFunction)
	}
	commandType := groups[0][1]
	subCommand := groups[0][2]

	var cue *Cue
	var err error
	switch commandType {
	case "set":
		cue, err = processSetCommand(subCommand)
	case "cycle":
		cue, err = processCycleCommand(subCommand)
	default:
		err = errors.New(commandErrorUndefinedFunction)
	}

	if err != nil {
		return nil, err
	}
	return cue, nil

}

func processCycleCommand(cmd string) (*Cue, error) {
	cue := Cue{}
	parts := strings.Split(cmd, ":")
	lightList := strings.Split(parts[0], ",")
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return nil, err
	}
	for x := range lightList {
		frame := Frame{}
		for y := 0; y < len(lightList); y++ {
			action := FrameAction{}
			action.LightName = lightList[y]

			action.NewState = light.State{
				RGB:      color.GetRGBFromString("#00FF00"),
				Duration: duration,
			}
			if x == y {
				action.NewState = light.State{
					RGB:      color.GetRGBFromString("#00FF00"),
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
			return nil, errors.New(commandErrorWrongPartCount)
		}
		if len(parts[0]) == 0 || len(parts[1]) == 0 || len(parts[2]) == 0 {
			return nil, errors.New(commandErrorWrongPartCount)
		}
		lightList := strings.Split(parts[0], ",")
		colorList := strings.Split(parts[1], ",")
		timeList := strings.Split(parts[2], ",")

		numTimes := len(timeList)
		numColors := len(colorList)
		numLights := len(lightList)

		if !(numTimes == numColors && numColors == numLights) {
			return nil, errors.New(commandErrorPartSizeMismatch)
		}
		for x := 0; x < numLights; x++ {
			action := FrameAction{}
			action.LightName = lightList[x]
			timeAsInt, err := strconv.Atoi(timeList[x])
			if err != nil {
				return nil, errors.New(commandErrorInvalidTime)
			}
			action.NewState = light.State{
				RGB:      color.GetRGBFromString(colorList[x]),
				Duration: time.Millisecond * time.Duration(timeAsInt),
			}
			frame.Actions = append(frame.Actions, action)
		}
		cue.Frames = append(cue.Frames, frame)
	}
	return &cue, nil
}
