package main

import (
	"fmt"

	"github.com/maodijim/tesla-api-go"
)

func main() {
	// login with refresh token
	// set global to false to login to China server
	teslaApi := tesla.NewTeslaApi("", "", "eyJ", true)
	// login with user password (will require browser pop up)
	teslaApi = tesla.NewTeslaApi("test@test.com", "test", "", true)
	fmt.Printf("Refresh Token: %s", teslaApi.RefreshToken())

	// required to get supercharging payment history
	// Can be found in the browser cookies "teslaSSORefreshToken" after login
	teslaApi.SetTeslaWebRefreshToken("eyJ")
	teslaApi.Login()

	// get vehicles in the account
	vehicles, _ := teslaApi.ListVehicles()
	// set active vehicle
	_ = teslaApi.SetActiveVehicle(vehicles[0])
	// get account supercharging payment history
	history, _ := teslaApi.GetSuperChargingHistory()
	fmt.Printf("Found %d supercharging history\n", len(history))
	// get active vehicle release notes
	note, _ := teslaApi.ReleaseNotes(false)
	fmt.Printf("Release note: %s\n", note)
	// set active vehicle cabin overheat protection temperature
	res, _ := teslaApi.SetCabinOverheatProtectionTemp(30)
	fmt.Println(res)
	// get active vehicle data
	_, _ = teslaApi.VehicleData()
	// get if active vehicle is in charging
	_ = teslaApi.IsCharging()
	// get if active vehicle is supercharging
	_ = teslaApi.IsFastCharging()
	// check if active vehicle has new software update
	version, hasUpdate := teslaApi.HasSoftwareUpdate()
	if hasUpdate {
		fmt.Printf("has software update %s\n", version)
	}
	// wake up active vehicle
	teslaApi.WakeUp()
	// unlock active vehicle door
	teslaApi.DoorUnlock()
	// lock active vehicle door
	teslaApi.DoorLock()
	// open active vehicle charge port
	teslaApi.ChargeDoorOpen()
	// close active vehicle charge port
	teslaApi.ChargeDoorClose()
	// set active vehicle to charge to max
	teslaApi.ChargeMaxRange()
	// set active vehicle to charge limit to x percent
	teslaApi.SetChargeLimit(50)

	// set active vehicle to climate temperature
	teslaApi.SetClimateTemp(23.5, 23.5)
	// set active vehicle to seat heater level
	teslaApi.SetSeatHeater(tesla.SeatFrontLeft, 3)
	// open active vehicle trunk (front or rear)
	teslaApi.ActuateTrunk(tesla.FrontTrunkType)
}
