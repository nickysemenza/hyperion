package homekit

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/nickysemenza/hyperion/control/trigger"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
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

//HomeKit represents and instance of a homekit manager
type HomeKit struct {
	master cue.MasterManager
}

func buildRawAccessoryList(accessoryList []Accessory) []*accessory.Accessory {
	accessories := make([]*accessory.Accessory, len(accessoryList))
	for i, a := range accessoryList {
		accessories[i] = a.RawAccessory
	}
	return accessories
}

func (hk *HomeKit) buildSwitchList(ctx context.Context) {
	//for now let's have N switches
	for x := 1; x <= numSwitches; x++ {
		id := x
		switchName := fmt.Sprintf("Switch %d", id)
		s := accessory.NewSwitch(accessory.Info{Name: switchName, Manufacturer: "hyperion"})
		s.Switch.On.OnValueRemoteUpdate(func(on bool) {
			if on {
				trigger.Action(ctx, "homekit-switch", id, hk.master)
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
func Start(ctx context.Context, wg *sync.WaitGroup, master cue.MasterManager) {
	defer wg.Done()
	hkConfig := config.GetServerConfig(ctx).Inputs.HomeKit
	if !hkConfig.Enabled {
		log.Info("homekit is not enabled")
		return
	}
	//config
	config := hc.Config{Pin: hkConfig.Pin, StoragePath: "./_homekit_data"}
	bridge := accessory.NewBridge(accessory.Info{Name: "bridge1", Manufacturer: "Hyperion"})

	//accessory setup
	hk := HomeKit{master: master}
	hk.buildSwitchList(ctx)
	accessories := buildRawAccessoryList(allAccessories)

	//start the server
	t, err := hc.NewIPTransport(config, bridge.Accessory, accessories...)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		<-ctx.Done()
		<-t.Stop()
		log.Println("Homekit services stopped.")
		return
	}()

	t.Start()
}
