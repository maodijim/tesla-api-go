package main

import "github.com/maodijim/tesla-api-go"

func main() {
	// login with refresh token
	teslaApi := tesla.NewTeslaApi("", "", "eyJ", true)
	// login with user password
	teslaApi = tesla.NewTeslaApi("test@test.com", "test", "", true)
	teslaApi.Login()
	vehicles, _ := teslaApi.ListVehicles()
	_ = teslaApi.SetActiveVehicle(vehicles[0])
	_, _ = teslaApi.VehicleData()
	teslaApi.WakeUp()
	teslaApi.DoorUnlock()
	teslaApi.DoorLock()
	teslaApi.ChargeDoorOpen()
	teslaApi.ChargeDoorClose()
	teslaApi.ChargeMaxRange()
	teslaApi.SetChargeLimit(50)
	teslaApi.SetClimateTemp(23.5, 23.5)
	teslaApi.SetSeatHeater(tesla.SeatFrontLeft, 3)
	teslaApi.ActuateTrunk(tesla.FrontTrunkType)
}
