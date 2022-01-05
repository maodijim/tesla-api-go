package tesla

import (
	"net/http"
	"strconv"
	"testing"
)

const (
	fakeId           = 249293153069797
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
			tt.ta.ListVehicleById(strconv.FormatInt(fakeVehId, 10))
			tt.ta.VehicleData()
			tt.ta.ChargeState()
			tt.ta.GetSuperChargingHistory()
			tt.ta.GetChargeHistory()
			tt.ta.SetClimateTemp(18.2, 18.2)
			tt.ta.NearByChargingSites()
			tt.ta.VehicleConfig()
			tt.ta.ChargeStandard()
		})
	}
}
