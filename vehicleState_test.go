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
			_, _ = tt.ta.GuiSetting()
			_, _ = tt.ta.MobileEnable()
			_, _ = tt.ta.VehicleConfig()
			_, _ = tt.ta.VehicleData()
			_, _ = tt.ta.VehicleState()
			_, _ = tt.ta.SoftwareUpdate()
		})
	}
}
