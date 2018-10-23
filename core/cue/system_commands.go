package cue

import "github.com/nickysemenza/hyperion/core/config"

const systemCommandBlackout = `
function process(timing)
	actions = {}
	for i, light in ipairs(light_list) do
	table.insert(actions, {
		light = light,
		color = "#000000",
		timing = timing
	})
	end
	frame1 = {
	actions = actions
	}
	return {
	frames = {frame1}
	}
end
`

var systemCommands = config.UserCommandMap{
	"blackout": config.UserCommand{Body: systemCommandBlackout},
}
