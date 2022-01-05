package tesla

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	cmdWakeUp       = "wake_up"
	cmdHonkHorn     = "honk_horn"
	cmdFlash        = "flash_lights"
	cmdHomeLink     = "trigger_homelink"
	cmdDoorUnlock   = "door_unlock"
	cmdActuateTrunk = "actuate_trunk"
	cmdSentryMode   = "set_sentry_mode"
	cmdWinsContr    = "window_control"
	cmdSunRoofContr = "sun_roof_control"
	cmdMediaToggle  = "media_toggle_playback"
	cmdMediaNext    = "media_next_track"
	cmdMediaPrev    = "media_prev_track"
	cmdMediaPrevFav = "media_prev_fav"
	cmdMediaNextFav = "media_next_fav"
	cmdMediaVolUp   = "media_volume_up"
	cmdMediaVolDown = "media_volume_down"
)

func (t *TeslaApi) WakeUp() (v *Vehicle, err error) {
	v = &Vehicle{}
	if t.activeVehicle.Id == 0 {
		return v, ErrNoActiveVehicle
	}
	u := joinPath(commandUrlBase, vehicleEndpoint, t.activeVehicle.GetIdStr(), cmdWakeUp)
	res, err := t.apiRequest(http.MethodPost, u, nil)
	if err != nil {
		return v, err
	}
	vRes, err := parseVehicleRes(res)
	if err != nil {
		return v, err
	}
	ve := vRes.Response.(Vehicle)
	v = &ve
	if res.StatusCode != 200 {
		r := BaseRes{}
		err = parseResp(res, &r)
		if err != nil {
			return v, errors.New(fmt.Sprintf("wake up return status code %d", res.StatusCode))
		}
		return v, errors.New(r.Err)
	}
	t.activeVehicle = *v
	return v, err
}

// HonkHorn Honks the horn twice.
func (t TeslaApi) HonkHorn() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdHonkHorn, "")
}

// FlashLights Flashes the headlights once
func (t TeslaApi) FlashLights() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdFlash, "")
}

// TriggerHomeLink Trigger homelink with current vehicle
func (t TeslaApi) TriggerHomeLink() (cmdRes *CommandsRes, err error) {
	ds, err := t.DriveState()
	if err != nil {
		return cmdRes, err
	}
	return t.sendCommand(cmdHomeLink, t.formUrlEncode(
		map[string]string{
			"lat": strconv.FormatFloat(ds.Latitude, 'f', -1, 64),
			"lon": strconv.FormatFloat(ds.Latitude, 'f', -1, 64),
		}))
}

func (t TeslaApi) SetSentryMod(on bool) (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdSentryMode, t.formUrlEncode(
		map[string]string{
			"on": strconv.FormatBool(on),
		}))
}

// DoorUnlock Unlocks the doors to the car. Extends the handles on the S and X.
func (t TeslaApi) DoorUnlock() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdDoorUnlock, "")
}

func (t TeslaApi) DoorLock() (cmdRes *CommandsRes, err error) {
	cmdRes, err = t.sendCommand(cmdDoorUnlock, "")
	return cmdRes, err
}

const (
	FrontTrunkType TrunkType = "front"
	RearTrunkType  TrunkType = "rear"
)

type TrunkType string

func (t TrunkType) String() string {
	return string(t)
}

func (t TeslaApi) ActuateTrunk(trunk TrunkType) (cmdRes *CommandsRes, err error) {
	cmdRes, err = t.sendCommand(cmdActuateTrunk, t.formUrlEncode(map[string]string{
		"which_trunk": trunk.String(),
	}))
	return cmdRes, err
}

const (
	WinCloseCmd WindowCmd = "close"
	WinVentCmd  WindowCmd = "vent"
)

type WindowCmd string

func (w WindowCmd) String() string {
	return string(w)
}

func (t TeslaApi) WindowControl(winCmd WindowCmd, lat, lon float64) (cmdRes *CommandsRes, err error) {
	cmdRes, err = t.sendCommand(cmdWinsContr, t.formUrlEncode(map[string]string{
		"command": winCmd.String(),
		"lat":     strconv.FormatFloat(lat, 'f', -1, 64),
		"lon":     strconv.FormatFloat(lon, 'f', -1, 64),
	}))
	return cmdRes, err
}

func (t TeslaApi) SunRoofControl(winCmd WindowCmd) (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdSunRoofContr, t.formUrlEncode(map[string]string{
		"state": winCmd.String(),
	}))
}

func (t TeslaApi) MediaToggle() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaToggle, "")
}

func (t TeslaApi) MediaNextTrack() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaNext, "")
}

func (t TeslaApi) MediaPrevTrack() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaPrev, "")
}

func (t TeslaApi) MediaPrevFav() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaPrevFav, "")
}

func (t TeslaApi) MediaNextFav() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaNextFav, "")
}

func (t TeslaApi) MediaVolUp() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaVolUp, "")
}

func (t TeslaApi) MediaVolDown() (cmdRes *CommandsRes, err error) {
	return t.sendCommand(cmdMediaVolDown, "")
}
