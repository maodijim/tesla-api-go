package tesla

import (
	"net/http"
	"reflect"
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
			_, _ = tt.ta.SetCabinOverheatProtectionTemp(40)
		})
	}
}

func TestTeslaApi_SetCabinOverheatProtectionTemp(t1 *testing.T) {
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
	type args struct {
		temp int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantCmdRes *CommandsRes
		wantErr    bool
	}{
		{
			name:    "Test set cabin overheat protection temp range",
			wantErr: true,
			args: args{
				temp: 0,
			},
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
			gotCmdRes, err := t.SetCabinOverheatProtectionTemp(tt.args.temp)
			if (err != nil) != tt.wantErr {
				t1.Errorf("SetCabinOverheatProtectionTemp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCmdRes, tt.wantCmdRes) {
				t1.Errorf("SetCabinOverheatProtectionTemp() gotCmdRes = %v, want %v", gotCmdRes, tt.wantCmdRes)
			}
		})
	}
}
