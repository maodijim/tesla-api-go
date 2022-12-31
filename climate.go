package tesla

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	cmdClimateAutoAcStart             = "auto_conditioning_start"
	cmdClimateAutoAcStop              = "auto_conditioning_stop"
	cmdClimateSetTemp                 = "set_temps"
	cmdClimateMax                     = "set_preconditioning_max"
	cmdClimateSeatHeater              = "remote_seat_heater_request"
	cmdClimateSteeringHeat            = "remote_steering_wheel_heater_request"
	cmdClimateCabinOverheatProtection = "set_cop_temp"
)

var (
	ErrInvalidSeatHeater   = errors.New("invalid heater must be between 0 to 5")
	ErrInvalidSeatHeaterLv = errors.New("invalid heater level must be between 0 to 3")
)

type ClimateState struct {
	BatteryHeater              bool    `json:"battery_heater"`
	BatteryHeaterNoPower       bool    `json:"battery_heater_no_power"`
	ClimateKeeperMode          string  `json:"climate_keeper_mode"`
	DefrostMode                int     `json:"defrost_mode"`
	DriverTempSetting          float64 `json:"driver_temp_setting"`
	FanStatus                  int     `json:"fan_status"`
	InsideTemp                 float64 `json:"inside_temp"`
	IsAutoConditioningOn       bool    `json:"is_auto_conditioning_on"`
	IsClimateOn                bool    `json:"is_climate_on"`
	IsFrontDefrosterOn         bool    `json:"is_front_defroster_on"`
	IsPreconditioning          bool    `json:"is_preconditioning"`
	IsRearDefrosterOn          bool    `json:"is_rear_defroster_on"`
	LeftTempDirection          int     `json:"left_temp_direction"`
	MaxAvailTemp               float64 `json:"max_avail_temp"`
	MinAvailTemp               float64 `json:"min_avail_temp"`
	OutsideTemp                float64 `json:"outside_temp"`
	PassengerTempSetting       float64 `json:"passenger_temp_setting"`
	RemoteHeaterControlEnabled bool    `json:"remote_heater_control_enabled"`
	RightTempDirection         int     `json:"right_temp_direction"`
	SeatHeaterLeft             int     `json:"seat_heater_left"`
	SeatHeaterRight            int     `json:"seat_heater_right"`
	SideMirrorHeaters          bool    `json:"side_mirror_heaters"`
	Timestamp                  int64   `json:"timestamp"`
	WiperBladeHeater           bool    `json:"wiper_blade_heater"`
}

func (t *TeslaApi) ClimateState() (cs *ClimateState, err error) {
	cs = &ClimateState{}
	lastUpdate := timestampSince(t.activeVehicleData.ClimateState.Timestamp)
	if lastUpdate < ClimateStateReqInterval && lastUpdate > 0 {
		return &t.activeVehicleData.ClimateState, nil
	}
	r := struct {
		BaseRes
		Response ClimateState `json:"response"`
	}{}
	err = t.sendDataRequest("climate_state", &r)
	if err != nil {
		return cs, err
	}
	cs = &r.Response
	t.activeVehicleData.ClimateState = *cs
	return cs, err
}

func (t *TeslaApi) ClimateAutoAcStart() (o *CommandsRes, err error) {
	return t.sendCommand(cmdClimateAutoAcStart, "")
}

func (t *TeslaApi) ClimateAutoAcStop() (o *CommandsRes, err error) {
	return t.sendCommand(cmdClimateAutoAcStop, "")
}

// SetClimateTemp temperature in celcius
func (t *TeslaApi) SetClimateTemp(driverTemp, passengerTemp float64) (o *CommandsRes, err error) {
	return t.sendCommand(
		cmdClimateSetTemp,
		t.jsonEncode(
			map[string]string{
				"driver_temp":    strconv.FormatFloat(driverTemp, 'f', 1, 64),
				"passenger_temp": strconv.FormatFloat(passengerTemp, 'f', 1, 64),
			},
		),
	)
}

func (t *TeslaApi) SetClimatePreConditionMax(on bool) (o *CommandsRes, err error) {
	return t.sendCommand(
		cmdClimateMax,
		t.jsonEncode(
			map[string]string{
				"on": strconv.FormatBool(on),
			},
		),
	)
}

const (
	SeatFrontLeft  SeatType = 0
	SeatFrontRight SeatType = 1
	SeatRearLeft   SeatType = 2
	SeatRearCenter SeatType = 4
	SeatRearRight  SeatType = 5
)

type SeatType int

// SetSeatHeater heater 0-5 , heat level 0-3
// heater seat
// 0 Front Left
// 1 Front right
// 2 Rear left
// 4 Rear center
// 5 Rear right
func (t *TeslaApi) SetSeatHeater(heater SeatType, level int) (o *CommandsRes, err error) {
	if heater < 0 || heater > 5 || heater == 3 {
		return o, ErrInvalidSeatHeater
	}
	if level < 0 || level > 3 {
		return o, ErrInvalidSeatHeaterLv
	}
	return t.sendCommand(
		cmdClimateSeatHeater,
		t.jsonEncode(
			map[string]string{
				"heater": strconv.FormatInt(int64(heater), 10),
				"level":  strconv.FormatInt(int64(level), 10),
			},
		),
	)
}

func (t *TeslaApi) SetSteeringHeater(on bool) (cmdRes *CommandsRes, err error) {
	return t.sendCommand(
		cmdClimateSteeringHeat,
		t.jsonEncode(
			map[string]string{
				"on": strconv.FormatBool(on),
			},
		),
	)
}

// SetCabinOverheatProtectionTemp set cabin overheat protection temperature in Celsius
func (t *TeslaApi) SetCabinOverheatProtectionTemp(temp int) (cmdRes *CommandsRes, err error) {
	copAllowedTemp := []int{30, 35, 40}
	if !contains(copAllowedTemp, temp) {
		return cmdRes, fmt.Errorf("invalid temp must be one of %v", copAllowedTemp)
	}
	return t.sendCommand(
		cmdClimateCabinOverheatProtection,
		t.jsonEncode(
			map[string]string{
				"cop_temp": strconv.FormatInt(int64(temp), 10),
			},
		),
	)
}
