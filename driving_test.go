package tesla

import (
	"testing"
)

func TestTeslaApi_RemoteStartDrive(t1 *testing.T) {
	tests := []struct {
		name string
		ta   *TeslaApi
	}{
		{
			name: "Test driving panics",
			ta:   fakeApi,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.ta.ActivateSpeedLimit("1234")
			tt.ta.ClearSpeedLimitPin("1234")
			tt.ta.DeactivateSpeedLimit("1234")
			tt.ta.DriveState()
			tt.ta.RemoteStartDrive()
			tt.ta.ResetValetPin()
			tt.ta.SetSpeedLimit(78)
			tt.ta.SetValetMode(true, "1234")
		})
	}
}
