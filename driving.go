package tesla

import (
	"errors"
	"strconv"
)

const (
	cmdDriveRemoteStart           = "remote_start_drive"
	cmdDriveSetLimit              = "speed_limit_set_limit"
	cmdDriverActivateSpeedLimit   = "speed_limit_activate"
	cmdDriverDeactivateSpeedLimit = "speed_limit_deactivate"
	cmdDriverClearSpeedLimitPin   = "speed_limit_clear_pin"
	cmdDriverSetValet             = "set_valet_mode"
	cmdDriveResetValetPin         = "reset_valet_pin"
)

var (
	ErrInvalidPinFormat = errors.New("invalid pin must be 4 digit pin")
)

type DriveState struct {
	GpsAsOf                 int     `json:"gps_as_of"`
	Heading                 int     `json:"heading"`
	Latitude                float64 `json:"latitude"`
	Longitude               float64 `json:"longitude"`
	NativeLatitude          float64 `json:"native_latitude"`
	NativeLocationSupported int     `json:"native_location_supported"`
	NativeLongitude         float64 `json:"native_longitude"`
	NativeType              string  `json:"native_type"`
	Power                   int     `json:"power"`
	ShiftState              string  `json:"shift_state"`
	Speed                   int     `json:"speed"`
	Timestamp               int64   `json:"timestamp"`
}

func (t *TeslaApi) DriveState() (ds *DriveState, err error) {
	ds = &DriveState{}
	lastUpdate := timestampSince(t.activeVehicleData.DriveState.Timestamp)
	if lastUpdate < DriveStateReqInterval && lastUpdate > 0 {
		return &t.activeVehicleData.DriveState, nil
	}
	r := struct {
		BaseRes
		Response DriveState `json:"response"`
	}{}
	err = t.sendDataRequest("drive_state", &r)
	if err != nil {
		return ds, err
	}
	ds = &r.Response
	t.activeVehicleData.DriveState = *ds
	return ds, err
}

func (t *TeslaApi) RemoteStartDrive() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdDriveRemoteStart, "")
}

// SetSpeedLimit speed limit between 50 - 90
func (t *TeslaApi) SetSpeedLimit(limitMph int) (cmdRes *CommandsRes, err error) {
	return t.sendCommand(
		cmdDriveSetLimit,
		t.formUrlEncode(
			map[string]string{
				"limit_mph": strconv.FormatInt(int64(limitMph), 10),
			},
		),
	)
}

// SetValetMode
// Valet Mode limits the car's top speed to 70MPH and 80kW of acceleration power.
// It also disables Homelink, Bluetooth and Wifi settings, and the ability to disable mobile access to the car.
// It also hides your favorites, home, and work locations in navigation.
//
// Note: the password parameter isn't required to turn on or off Valet Mode, even with a previous PIN set.
// If you clear the PIN and activate Valet Mode without the parameter,
//  you will only be able to deactivate it from your car's screen by signing into your Tesla account.
func (t *TeslaApi) SetValetMode(on bool, password string) (cmdRes *CommandsRes, err error) {
	return t.sendCommand(
		cmdDriverSetValet,
		t.formUrlEncode(
			map[string]string{
				"on":       strconv.FormatBool(on),
				"password": password,
			},
		),
	)
}

func (t *TeslaApi) ResetValetPin() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(
		cmdDriveResetValetPin,
		"",
	)
}

func (t *TeslaApi) ActivateSpeedLimit(pin string) (cmdRes *CommandsRes, err error) {
	if len(pin) != 4 {
		return cmdRes, ErrInvalidPinFormat
	}
	return t.sendCommand(cmdDriverActivateSpeedLimit, t.formUrlEncode(
		map[string]string{
			"pin": pin,
		},
	))
}

func (t *TeslaApi) DeactivateSpeedLimit(pin string) (cmdRes *CommandsRes, err error) {
	if len(pin) != 4 {
		return cmdRes, ErrInvalidPinFormat
	}
	return t.sendCommand(
		cmdDriverDeactivateSpeedLimit,
		t.formUrlEncode(
			map[string]string{
				"pin": pin,
			},
		),
	)
}

func (t *TeslaApi) ClearSpeedLimitPin(pin string) (cmdRes *CommandsRes, err error) {
	if len(pin) != 4 {
		return cmdRes, ErrInvalidPinFormat
	}
	return t.sendCommand(
		cmdDriverClearSpeedLimitPin,
		t.formUrlEncode(
			map[string]string{
				"pin": pin,
			},
		),
	)
}
