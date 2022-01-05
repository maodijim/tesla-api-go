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
			tt.ta.ActuateTrunk(FrontTrunkType)
			tt.ta.DoorLock()
			tt.ta.DoorLock()
			tt.ta.FlashLights()
			tt.ta.HonkHorn()
			tt.ta.MediaNextTrack()
			tt.ta.MediaPrevTrack()
			tt.ta.MediaNextFav()
			tt.ta.MediaPrevFav()
			tt.ta.MediaToggle()
			tt.ta.MediaVolUp()
			tt.ta.MediaVolDown()
			tt.ta.SetSentryMod(true)
			tt.ta.SunRoofControl(WinVentCmd)
			tt.ta.TriggerHomeLink()
			tt.ta.WakeUp()
			tt.ta.WindowControl(WinCloseCmd, 100.00, 100.00)

		})
	}
}
