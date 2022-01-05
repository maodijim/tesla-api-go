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
			_, _ = tt.ta.ActivateSpeedLimit("1234")
			_, _ = tt.ta.ClearSpeedLimitPin("1234")
			_, _ = tt.ta.DeactivateSpeedLimit("1234")
			_, _ = tt.ta.DriveState()
			_, _ = tt.ta.RemoteStartDrive()
			_, _ = tt.ta.ResetValetPin()
			_, _ = tt.ta.SetSpeedLimit(78)
			_, _ = tt.ta.SetValetMode(true, "1234")
		})
	}
}
