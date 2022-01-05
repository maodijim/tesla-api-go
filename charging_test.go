package tesla

import (
	"testing"
)

func TestTeslaApi_Charging(t1 *testing.T) {
	tests := []struct {
		name string
		ta   *TeslaApi
	}{
		{
			name: "Test charging panics",
			ta:   fakeApi,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			_, _ = tt.ta.ChargeState()
			_, _ = tt.ta.ChargeDoorOpen()
			_, _ = tt.ta.ChargeDoorClose()
			_, _ = tt.ta.ChargeMaxRange()
			_, _ = tt.ta.ChargeStandard()
			_, _ = tt.ta.ChargeStart()
			_, _ = tt.ta.ChargeStop()
			_, _ = tt.ta.NearByChargingSites()
			_, _ = tt.ta.GetSuperChargingHistory()
			_, _ = tt.ta.GetChargeHistory()
			_, _ = tt.ta.SetScheduledDeparture(true, 100, 100, true, true, true, true)
			_, _ = tt.ta.SetScheduledCharge(true, 120)
			tt.ta.IsCharging()
		})
	}
}
