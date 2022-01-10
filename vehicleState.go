package tesla

import (
	"errors"
	"net/http"
)

type GuiSettings struct {
	Gui24HourTime       bool   `json:"gui_24_hour_time"`
	GuiChargeRateUnits  string `json:"gui_charge_rate_units"`
	GuiDistanceUnits    string `json:"gui_distance_units"`
	GuiRangeDisplay     string `json:"gui_range_display"`
	GuiTemperatureUnits string `json:"gui_temperature_units"`
	ShowRangeUnits      bool   `json:"show_range_units"`
	Timestamp           int64  `json:"timestamp"`
}

type VehicleState struct {
	ApiVersion          int    `json:"api_version"`
	AutoparkStateV2     string `json:"autopark_state_v2"`
	AutoparkStyle       string `json:"autopark_style"`
	CalendarSupported   bool   `json:"calendar_supported"`
	CarVersion          string `json:"car_version"`
	CenterDisplayState  int    `json:"center_display_state"`
	Df                  int    `json:"df"`
	Dr                  int    `json:"dr"`
	FdWindow            int    `json:"fd_window"`
	FpWindow            int    `json:"fp_window"`
	Ft                  int    `json:"ft"`
	HomelinkDeviceCount int    `json:"homelink_device_count"`
	HomelinkNearby      bool   `json:"homelink_nearby"`
	IsUserPresent       bool   `json:"is_user_present"`
	LastAutoparkError   string `json:"last_autopark_error"`
	Locked              bool   `json:"locked"`
	MediaState          struct {
		RemoteControlEnabled bool `json:"remote_control_enabled"`
	} `json:"media_state"`
	NotificationsSupported  bool    `json:"notifications_supported"`
	Odometer                float64 `json:"odometer"`
	ParsedCalendarSupported bool    `json:"parsed_calendar_supported"`
	Pf                      int     `json:"pf"`
	Pr                      int     `json:"pr"`
	RdWindow                int     `json:"rd_window"`
	RemoteStart             bool    `json:"remote_start"`
	RemoteStartEnabled      bool    `json:"remote_start_enabled"`
	RemoteStartSupported    bool    `json:"remote_start_supported"`
	RpWindow                int     `json:"rp_window"`
	Rt                      int     `json:"rt"`
	SentryMode              bool    `json:"sentry_mode"`
	SentryModeAvailable     bool    `json:"sentry_mode_available"`
	SmartSummonAvailable    bool    `json:"smart_summon_available"`
	SoftwareUpdate          struct {
		DownloadPerc        int    `json:"download_perc"`
		ExpectedDurationSec int    `json:"expected_duration_sec"`
		InstallPerc         int    `json:"install_perc"`
		Status              string `json:"status"`
		Version             string `json:"version"`
	} `json:"software_update"`
	SpeedLimitMode struct {
		Active          bool    `json:"active"`
		CurrentLimitMph float64 `json:"current_limit_mph"`
		MaxLimitMph     float64 `json:"max_limit_mph"`
		MinLimitMph     float64 `json:"min_limit_mph"`
		PinCodeSet      bool    `json:"pin_code_set"`
	} `json:"speed_limit_mode"`
	SummonStandbyModeEnabled bool        `json:"summon_standby_mode_enabled"`
	SunRoofPercentOpen       int         `json:"sun_roof_percent_open"`
	SunRoofState             string      `json:"sun_roof_state"`
	Timestamp                int64       `json:"timestamp"`
	ValetMode                bool        `json:"valet_mode"`
	ValetPinNeeded           bool        `json:"valet_pin_needed"`
	VehicleName              interface{} `json:"vehicle_name"`
	VehicleSelfTestProgress  int         `json:"vehicle_self_test_progress"`
	VehicleSelfTestRequested bool        `json:"vehicle_self_test_requested"`
	WebcamAvailable          bool        `json:"webcam_available"`
}

type VehicleConfig struct {
	CanAcceptNavigationRequests bool   `json:"can_accept_navigation_requests"`
	CanActuateTrunks            bool   `json:"can_actuate_trunks"`
	CarSpecialType              string `json:"car_special_type"`
	CarType                     string `json:"car_type"`
	ChargePortType              string `json:"charge_port_type"`
	DefaultChargeToMax          bool   `json:"default_charge_to_max"`
	EceRestrictions             bool   `json:"ece_restrictions"`
	EuVehicle                   bool   `json:"eu_vehicle"`
	ExteriorColor               string `json:"exterior_color"`
	HasAirSuspension            bool   `json:"has_air_suspension"`
	HasLudicrousMode            bool   `json:"has_ludicrous_mode"`
	MotorizedChargePort         bool   `json:"motorized_charge_port"`
	Plg                         bool   `json:"plg"`
	RearSeatHeaters             int    `json:"rear_seat_heaters"`
	RearSeatType                int    `json:"rear_seat_type"`
	Rhd                         bool   `json:"rhd"`
	RoofColor                   string `json:"roof_color"`
	SeatType                    int    `json:"seat_type"`
	SpoilerType                 string `json:"spoiler_type"`
	SunRoofInstalled            int    `json:"sun_roof_installed"`
	ThirdRowSeats               string `json:"third_row_seats"`
	Timestamp                   int64  `json:"timestamp"`
	TrimBadging                 string `json:"trim_badging"`
	UseRangeBadging             bool   `json:"use_range_badging"`
	WheelType                   string `json:"wheel_type"`
}

type VehicleData struct {
	Vehicle
	DriveState    DriveState    `json:"drive_state"`
	ClimateState  ClimateState  `json:"climate_state"`
	ChargeState   ChargeState   `json:"charge_state"`
	GuiSettings   GuiSettings   `json:"gui_settings"`
	VehicleState  VehicleState  `json:"vehicle_state"`
	VehicleConfig VehicleConfig `json:"vehicle_config"`
}

type VehicleDataRes struct {
	BaseRes
	Response VehicleData `json:"response"`
}

func (t *TeslaApi) VehicleData() (vd *VehicleData, err error) {
	vd = &VehicleData{}
	if t.activeVehicle.Id == 0 {
		return vd, ErrNoActiveVehicle
	}
	u := joinPath(commandUrlBase, vehicleEndpoint, t.activeVehicle.GetIdStr(), "vehicle_data")
	res, err := t.apiRequest(http.MethodGet, u, nil)
	if err != nil {
		return vd, err
	}
	vdRes := VehicleDataRes{}
	err = parseResp(res, &vdRes)
	if err != nil {
		return vd, errors.New(vdRes.Err)
	}
	vd = &vdRes.Response
	t.activeVehicleData = *vd
	return vd, err
}

func (t *TeslaApi) GuiSetting() (gs *GuiSettings, err error) {
	gs = &GuiSettings{}
	lastUpdate := timestampSince(t.activeVehicleData.GuiSettings.Timestamp)
	if lastUpdate < GuiSettingReqInterval && lastUpdate > 0 {
		return &t.activeVehicleData.GuiSettings, nil
	}
	r := struct {
		BaseRes
		Response GuiSettings `json:"response"`
	}{}
	err = t.sendDataRequest("gui_settings", &r)
	if err != nil {
		return gs, err
	}
	gs = &r.Response
	t.activeVehicleData.GuiSettings = *gs
	return gs, err
}

func (t *TeslaApi) VehicleState() (vs *VehicleState, err error) {
	vs = &VehicleState{}
	lastUpdate := timestampSince(t.activeVehicleData.VehicleState.Timestamp)
	if lastUpdate < VehicleStateReqInterval && lastUpdate > 0 {
		return &t.activeVehicleData.VehicleState, nil
	}
	r := struct {
		BaseRes
		Response VehicleState `json:"response"`
	}{}
	err = t.sendDataRequest("vehicle_state", &r)
	if err != nil {
		return vs, err
	}
	vs = &r.Response
	t.activeVehicleData.VehicleState = *vs
	return vs, err
}

func (t *TeslaApi) VehicleConfig() (vc *VehicleConfig, err error) {
	vc = &VehicleConfig{}
	lastUpdate := timestampSince(t.activeVehicleData.VehicleConfig.Timestamp)
	if lastUpdate < VehicleConfigReqInterval && lastUpdate > 0 {
		return &t.activeVehicleData.VehicleConfig, nil
	}
	r := struct {
		BaseRes
		Response VehicleConfig `json:"response"`
	}{}
	err = t.sendDataRequest("vehicle_config", &r)
	if err != nil {
		return vc, err
	}
	vc = &r.Response
	t.activeVehicleData.VehicleConfig = *vc
	return vc, err
}

func (t TeslaApi) MobileEnable() (isEnabled bool, err error) {
	r := struct {
		BaseRes
		Response bool `json:"response"`
	}{}
	err = t.sendDataRequest("mobile_enabled", &r)
	if err != nil {
		return isEnabled, err
	}
	return r.Response, err
}
