package tesla

import (
	"testing"
)

func TestTeslaApi_Climate(t1 *testing.T) {
	tests := []struct {
		name string
		ta   *TeslaApi
	}{
		{
			name: "Test climate panics",
			ta:   fakeApi,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.ta.ClimateAutoAcStart()
			tt.ta.ClimateAutoAcStop()
			tt.ta.ClimateState()
			tt.ta.SetClimatePreConditionMax(true)
			tt.ta.SetClimateTemp(23.5, 23.5)
			tt.ta.SetSeatHeater(SeatFrontLeft, 3)
			tt.ta.SetSteeringHeater(true)
		})
	}
}
