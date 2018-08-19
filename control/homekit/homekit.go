package homekit

import (
	"fmt"
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/nickysemenza/hyperion/control/trigger"
)

const numSwitches = 6

type accessoryType string

const (
	button accessoryType = "button"
)

var allAccessories []Accessory

//Accessory wraps a hc/accessory
type Accessory struct {
	RawAccessory *accessory.Accessory
	Type         accessoryType
	Name         string
}

func buildRawAccessoryList(accessoryList []Accessory) []*accessory.Accessory {
	accessories := make([]*accessory.Accessory, len(accessoryList))
	for i, a := range accessoryList {
		accessories[i] = a.RawAccessory
	}
	return accessories
}

func buildSwitchList() {
	//for now let's have N switches
	for x := 1; x <= numSwitches; x++ {
		id := x
		switchName := fmt.Sprintf("Switch %d", id)
		s := accessory.NewSwitch(accessory.Info{Name: switchName, Manufacturer: "hyperion"})
		s.Switch.On.OnValueRemoteUpdate(func(on bool) {
			if on {
				trigger.Action("homekit-switch", id)
				s.Switch.On.SetValue(false)
			}
			log.Printf("[homekit] changed: [%s] to %t", s.Accessory.Info.Name.String.GetValue(), on)
		})
		allAccessories = append(allAccessories, Accessory{
			RawAccessory: s.Accessory,
			Type:         button,
			Name:         switchName,
		})
	}
}

//Start starts the HomeKit services
func Start() {
	//config
	config := hc.Config{Pin: "10000000", StoragePath: "./_homekit_data"}
	bridge := accessory.NewBridge(accessory.Info{Name: "bridge1", Manufacturer: "Hyperion"})

	//accessory setup
	buildSwitchList()
	accessories := buildRawAccessoryList(allAccessories)

	//start the server
	t, err := hc.NewIPTransport(config, bridge.Accessory, accessories...)
	if err != nil {
		log.Panic(err)
	}

	//shutdown handler
	hc.OnTermination(func() {
		<-t.Stop()
		log.Println("Homekit services stopped.")
	})

	t.Start()
}
