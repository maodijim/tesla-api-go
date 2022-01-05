package tesla

import (
	"testing"
)

func TestTeslaApi_VehicleState(t1 *testing.T) {
	tests := []struct {
		name string
		ta   *TeslaApi
	}{
		{
			name: "Test vehicle data panics",
			ta:   fakeApi,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.ta.GuiSetting()
			tt.ta.MobileEnable()
			tt.ta.VehicleConfig()
			tt.ta.VehicleData()
			tt.ta.VehicleState()
		})
	}
}
