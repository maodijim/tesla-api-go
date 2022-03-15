package tesla

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

const (
	fakeId           = 36944823769655690
	fakeIdS          = "36944823769655690"
	fakeVehId        = 3165661272
	fakeAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ilg0RmNua0RCUVBUTnBrZTZiMnNuRi04YmdVUSJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92MyIsImF1ZCI6WyJodHRwczovL293bmVyLWFwaS50ZXNsYW1vdG9ycy5jb20vIiwiaHR0cHM6Ly9hdXRoLnRlc2xhLmNvbS9vYXV0aDIvdjMvdXNlcmluZm8iXSwiYXpwIjoib3duZXJhcGkiLCJzdWIiOiIxMjMyMzIzMjU2Ny0xMjM0LTEyMzQtYjQ1Mi02MTFmOTYiLCJzY3AiOlsib3BlbmlkIiwiZW1haWwiLCJvZmZsaW5lX2FjY2VzcyJdLCJhbXIiOlsicHdkIl0sImV4cCI6MTY0MTM3NTQyOCwiaWF0IjoxNjQxMzQ2NjI4LCJhdXRoX3RpbWUiOjE2NDEwMTE1NzJ9.pgmUIpbfTFseIl7top08x7LVk4Bmv06MyKUXJqRweuY"
	fakeRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ilg0RmNua0RCUVBUTnBrZTZiMnNuRi04YmdVUSJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY29tL29hdXRoMi92MyIsImF1ZCI6Imh0dHBzOi8vYXV0aC50ZXNsYS5jb20vb2F1dGgyL3YzL3Rva2VuIiwiaWF0IjoxNjQxMzM4MDg1LCJzY3AiOlsib3BlbmlkIiwib2ZmbGluZV9hY2Nlc3MiXSwiZGF0YSI6eyJ2IjoiMSIsImF1ZCI6Imh0dHBzOi8vb3duZXItYXBpLnRlc2xhbW90b3JzLmNvbS8iLCJzdWIiOiI3MTYxMjM0NTY2Ny0xMjM0LTEyMzQtMTIzNC1naGhqajU0MjIzMjE0YSIsInNjcCI6WyJvcGVuaWQiLCJlbWFpbCIsIm9mZmxpbmVfYWNjZXNzIl0sImF6cCI6Im93bmVyYXBpIiwiYW1yIjpbInB3ZCIsIm1mYSIsIm90cCJdLCJhdXRoX3RpbWUiOjE2NDEzMzc5OTN9fQ.LqhwuZCtZibRQbTn-OwA_oAOsKIMIHYCPPXdpYz_ZhQ"
)

var (
	fakeVehicle = Vehicle{
		Id:        fakeId,
		VehicleId: fakeVehId,
	}
	fakeApi = &TeslaApi{
		activeVehicle: Vehicle{Id: fakeId},
		client:        &http.Client{},
		accessToken:   fakeAccessToken,
		refreshToken:  fakeRefreshToken,
	}
)

func TestTeslaApi_Vehicles(t1 *testing.T) {
	tests := []struct {
		name   string
		ta     *TeslaApi
		wantVs []Vehicle
	}{
		{
			name: "Test Vehicle panics",
			ta:   fakeApi,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			_, _ = tt.ta.ListVehicleById(strconv.FormatInt(fakeVehId, 10))
			_, _ = tt.ta.VehicleData()
			_, _ = tt.ta.ChargeState()
			_, _ = tt.ta.GetSuperChargingHistory()
			_, _ = tt.ta.GetChargeHistory()
			_, _ = tt.ta.SetClimateTemp(18.2, 18.2)
			_, _ = tt.ta.NearByChargingSites()
			_, _ = tt.ta.VehicleConfig()
			_, _ = tt.ta.ChargeStandard()
			_, _ = tt.ta.HasSoftwareUpdate()
			_, _ = tt.ta.IsSoftwareInstalling()
		})
	}
}

func Test_convertMapToVehicle(t *testing.T) {
	type args struct {
		in map[string]interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantVehicle Vehicle
	}{
		{
			name: "Test Vehicle Struct Map conversion id int type",
			args: args{
				in: map[string]interface{}{
					"id": fakeId,
				},
			},
			wantVehicle: Vehicle{
				Id: fakeId,
			},
		},
		{
			name: "Test Vehicle Struct Map conversion id json.Number type",
			args: args{
				in: map[string]interface{}{
					"id": json.Number(fakeIdS),
				},
			},
			wantVehicle: Vehicle{
				Id: fakeId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVehicle := convertMapToVehicle(tt.args.in); !reflect.DeepEqual(gotVehicle, tt.wantVehicle) {
				t.Errorf("convertMapToVehicle() = %v, want %v", gotVehicle, tt.wantVehicle)
			}
		})
	}
}
