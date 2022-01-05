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
			tt.ta.ChargeState()
			tt.ta.ChargeDoorOpen()
			tt.ta.ChargeDoorClose()
			tt.ta.ChargeMaxRange()
			tt.ta.ChargeStandard()
			tt.ta.ChargeStart()
			tt.ta.ChargeStop()
			tt.ta.NearByChargingSites()
			tt.ta.GetSuperChargingHistory()
			tt.ta.GetChargeHistory()
			tt.ta.SetScheduledDeparture(true, 100, 100, true, true, true, true)
			tt.ta.SetScheduledCharge(true, 120)
			tt.ta.IsCharging()
		})
	}
}
