package homekit

import (
	"reflect"
	"testing"

	"github.com/brutella/hc/accessory"
)

func Test_buildRawAccessoryList(t *testing.T) {
	a1 := accessory.NewSwitch(accessory.Info{}).Accessory
	tests := []struct {
		name          string
		accessoryList []Accessory
		want          []*accessory.Accessory
	}{
		{"empty", []Accessory{}, []*accessory.Accessory{}},
		{"simple", []Accessory{Accessory{RawAccessory: a1}}, []*accessory.Accessory{a1}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildRawAccessoryList(tt.accessoryList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildRawAccessoryList() = %v, want %v", got, tt.want)
			}
		})
	}
}
