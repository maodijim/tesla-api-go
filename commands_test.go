package tesla

import (
	"testing"
)

func TestTeslaApi_Command(t1 *testing.T) {
	tests := []struct {
		name string
		ta   *TeslaApi
	}{
		{
			name: "Test command panics",
			ta:   fakeApi,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			_, _ = tt.ta.ActuateTrunk(FrontTrunkType)
			_, _ = tt.ta.DoorLock()
			_, _ = tt.ta.DoorLock()
			_, _ = tt.ta.FlashLights()
			_, _ = tt.ta.HonkHorn()
			_, _ = tt.ta.MediaNextTrack()
			_, _ = tt.ta.MediaPrevTrack()
			_, _ = tt.ta.MediaNextFav()
			_, _ = tt.ta.MediaPrevFav()
			_, _ = tt.ta.MediaToggle()
			_, _ = tt.ta.MediaVolUp()
			_, _ = tt.ta.MediaVolDown()
			_, _ = tt.ta.SetSentryMode(true)
			_, _ = tt.ta.SunRoofControl(WinVentCmd)
			_, _ = tt.ta.TriggerHomeLink()
			_, _ = tt.ta.WakeUp()
			_, _ = tt.ta.WindowControl(WinCloseCmd, 100.00, 100.00)
			_, _ = tt.ta.ScheduleSoftwareUpdate(100)
			_, _ = tt.ta.CancelSoftwareUpdate()
		})
	}
}
