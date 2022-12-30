package tesla

import (
	"net/http"
	"reflect"
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

func TestTeslaApi_RemoteBoomBox(t1 *testing.T) {
	type fields struct {
		AuthReq           AuthReq
		client            *http.Client
		cookies           []*http.Cookie
		activeVehicle     Vehicle
		activeVehicleData VehicleData
		accessToken       string
		refreshToken      string
		renewingToken     bool
	}
	tests := []struct {
		name       string
		fields     fields
		wantCmdRes *CommandsRes
		wantErr    bool
	}{
		{
			name: "Test remote boom box version check",
			fields: fields{
				activeVehicleData: VehicleData{
					VehicleState: VehicleState{
						CarVersion: "2022.43.25.1",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TeslaApi{
				AuthReq:           tt.fields.AuthReq,
				client:            tt.fields.client,
				cookies:           tt.fields.cookies,
				activeVehicle:     tt.fields.activeVehicle,
				activeVehicleData: tt.fields.activeVehicleData,
				accessToken:       tt.fields.accessToken,
				refreshToken:      tt.fields.refreshToken,
				renewingToken:     tt.fields.renewingToken,
			}
			gotCmdRes, err := t.RemoteBoomBox()
			if (err != nil) != tt.wantErr {
				t1.Errorf("RemoteBoomBox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCmdRes, tt.wantCmdRes) {
				t1.Errorf("RemoteBoomBox() gotCmdRes = %v, want %v", gotCmdRes, tt.wantCmdRes)
			}
		})
	}
}
