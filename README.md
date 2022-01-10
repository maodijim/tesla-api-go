# tesla-api-go
![Test](https://github.com/maodijim/tesla-api-go/actions/workflows/go.yml/badge.svg)

This is an unofficial Go Tesla API client based on the documentation https://github.com/timdorr/tesla-api

## Installation
```sh
go get github.com/maodijim/tesla-api-go
```

## Usage
```go
import "github.com/maodijim/tesla-api-go"
// with credential
teslaApi := tesla.NewTeslaApi(username, password, "", true)

// with refresh token
teslaApi := tesla.NewTeslaApi("", "", "eyJ...", true)

teslaApi.Login()
vehicles, err := teslaApi.ListVehicles()

teslaApi.SetActiveVehicle(vehicles[0])

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

```


## Implemented
- Auth
  - [x] username & password 
  - [x] refresh token
  - [x] MFA
- Vehicle
- [x] /api/1/vehicles
- [x] /api/1/vehicles/{id}
- State
  - [x] /api/1/vehicles/{id}/vehicle_data
  - [x] /api/1/vehicles/{id}/data_request/charge_state
  - [x] /api/1/vehicles/{id}/data_request/climate_state
  - [x] /api/1/vehicles/{id}/data_request/drive_state
  - [x] /api/1/vehicles/{id}/data_request/gui_settings
  - [x] /api/1/vehicles/{id}/data_request/vehicle_state
  - [x] /api/1/vehicles/{id}/data_request/vehicle_config
  - [x] /api/1/vehicles/{id}/mobile_enabled
  - [x] /api/1/vehicles/{id}/nearby_charging_sites
- Commands
  - Wake
    - [x] /api/1/vehicles/{id}/wake_up
  - Alerts
    - [x] /api/1/vehicles/{id}/command/honk_horn
    - [x] /api/1/vehicles/{id}/command/flash_lights
  - Remote Start
    - [x] /api/1/vehicles/{id}/command/remote_start_drive
  - Homelink
    - [x] /api/1/vehicles/{id}/command/trigger_homelink
  - /api/1/vehicles/{id}/command/speed_limit_set_limit
    - [x] /api/1/vehicles/{id}/command/speed_limit_set_limit
    - [x] /api/1/vehicles/{id}/command/speed_limit_activate
    - [x] /api/1/vehicles/{id}/command/speed_limit_deactivate
    - [x] /api/1/vehicles/{id}/command/speed_limit_clear_pin
  - Valet Mode
    - [x] /api/1/vehicles/{id}/command/set_valet_mode
    - [x] /api/1/vehicles/{id}/command/reset_valet_pin
  - Sentry Mode
    - [x] /api/1/vehicles/{id}/command/set_sentry_mode
  - Doors
    - [x] /api/1/vehicles/{id}/command/door_unlock
    - [x] /api/1/vehicles/{id}/command/door_lock
  - Trunk
    - [x] /api/1/vehicles/{id}/command/actuate_trunk
  - Windows
    - [x] /api/1/vehicles/{id}/command/window_control
  - Sunroof
    - [x] /api/1/vehicles/{id}/command/sun_roof_control
  - Charging
    - [x] /api/1/vehicles/{id}/command/charge_port_door_open
    - [x] /api/1/vehicles/{id}/command/charge_port_door_close
    - [x] /api/1/vehicles/{id}/command/charge_start
    - [x] /api/1/vehicles/{id}/command/charge_stop
    - [x] /api/1/vehicles/{id}/command/charge_standard
    - [x] /api/1/vehicles/{id}/command/charge_max_range
    - [x] /api/1/vehicles/{id}/command/set_charge_limit
    - [x] /api/1/vehicles/{id}/command/set_charging_amps
    - [x] /api/1/vehicles/{id}/command/set_scheduled_charging
    - [x] /api/1/vehicles/{id}/command/set_scheduled_departure
    - [x] api/1/vehicles/{vehicle_id}/charge_history
    - [x] super charger history
  - Climate
    - [x] /api/1/vehicles/{id}/command/auto_conditioning_start
    - [x] /api/1/vehicles/{id}/command/auto_conditioning_stop
    - [x] /api/1/vehicles/{id}/command/set_temps
    - [x] /api/1/vehicles/{id}/command/set_preconditioning_max
    - [x] /api/1/vehicles/{id}/command/remote_seat_heater_request
    - [x] /api/1/vehicles/{id}/command/remote_steering_wheel_heater_request
  - Media
    - [x] /api/1/vehicles/{id}/command/media_toggle_playback
    - [x] /api/1/vehicles/{id}/command/media_next_track
    - [x] /api/1/vehicles/{id}/command/media_prev_track
    - [x] /api/1/vehicles/{id}/command/media_next_fav
    - [x] /api/1/vehicles/{id}/command/media_prev_fav
    - [x] /api/1/vehicles/{id}/command/media_volume_up
    - [x] /api/1/vehicles/{id}/command/media_volume_down
  - Sharing
    - [ ] /api/1/vehicles/{id}/command/share
  - Software Updates
    - [x] /api/1/vehicles/{id}/command/schedule_software_update
    - [x] /api/1/vehicles/{id}/command/cancel_software_update
- [ ] Streaming
- [ ] Autopark/Summon
- [ ] Solar
- [ ] Powerwall
