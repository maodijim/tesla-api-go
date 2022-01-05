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
			_, _ = tt.ta.ClimateAutoAcStart()
			_, _ = tt.ta.ClimateAutoAcStop()
			_, _ = tt.ta.ClimateState()
			_, _ = tt.ta.SetClimatePreConditionMax(true)
			_, _ = tt.ta.SetClimateTemp(23.5, 23.5)
			_, _ = tt.ta.SetSeatHeater(SeatFrontLeft, 3)
			_, _ = tt.ta.SetSteeringHeater(true)
		})
	}
}
