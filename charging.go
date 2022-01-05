package tesla

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	cmdChargeDoorOpen  = "charge_port_door_open"
	cmdChargeDoorClose = "charge_port_door_close"
	cmdChargeStart     = "charge_start"
	cmdChargeStop      = "charge_stop"
	cmdChargeStandard  = "charge_standard"
	cmdChargeMax       = "charge_max_range"
	cmdChargeLimit     = "set_charge_limit"
	cmdChargeAmp       = "set_charging_amps"
	cmdChargeSchedule  = "set_scheduled_charging"
	cmdChargeDeparture = "set_scheduled_departure"
)

var (
	ErrInvalidAmps          = errors.New("invalid charge amps must be within 1 to 1000")
	ErrInvalidChargePercent = errors.New("invalid charge limit must be within 0 to 100")
	ErrInvalidChargeTime    = errors.New("invalid charge time must be within 0 to 1440")
)

type ChargeState struct {
	BatteryHeaterOn             bool        `json:"battery_heater_on"`
	BatteryLevel                int         `json:"battery_level"`
	BatteryRange                float64     `json:"battery_range"`
	ChargeCurrentRequest        int         `json:"charge_current_request"`
	ChargeCurrentRequestMax     int         `json:"charge_current_request_max"`
	ChargeEnableRequest         bool        `json:"charge_enable_request"`
	ChargeEnergyAdded           float64     `json:"charge_energy_added"`
	ChargeLimitSoc              int         `json:"charge_limit_soc"`
	ChargeLimitSocMax           int         `json:"charge_limit_soc_max"`
	ChargeLimitSocMin           int         `json:"charge_limit_soc_min"`
	ChargeLimitSocStd           int         `json:"charge_limit_soc_std"`
	ChargeMilesAddedIdeal       float64     `json:"charge_miles_added_ideal"`
	ChargeMilesAddedRated       float64     `json:"charge_miles_added_rated"`
	ChargePortColdWeatherMode   interface{} `json:"charge_port_cold_weather_mode"`
	ChargePortDoorOpen          bool        `json:"charge_port_door_open"`
	ChargePortLatch             string      `json:"charge_port_latch"`
	ChargeRate                  float64     `json:"charge_rate"`
	ChargeToMaxRange            bool        `json:"charge_to_max_range"`
	ChargerActualCurrent        int         `json:"charger_actual_current"`
	ChargerPhases               int         `json:"charger_phases"`
	ChargerPilotCurrent         int         `json:"charger_pilot_current"`
	ChargerPower                int         `json:"charger_power"`
	ChargerVoltage              int         `json:"charger_voltage"`
	ChargingState               string      `json:"charging_state"`
	ConnChargeCable             string      `json:"conn_charge_cable"`
	EstBatteryRange             float64     `json:"est_battery_range"`
	FastChargerBrand            string      `json:"fast_charger_brand"`
	FastChargerPresent          bool        `json:"fast_charger_present"`
	FastChargerType             string      `json:"fast_charger_type"`
	IdealBatteryRange           float64     `json:"ideal_battery_range"`
	ManagedChargingActive       bool        `json:"managed_charging_active"`
	ManagedChargingStartTime    int         `json:"managed_charging_start_time"`
	ManagedChargingUserCanceled bool        `json:"managed_charging_user_canceled"`
	MaxRangeChargeCounter       int         `json:"max_range_charge_counter"`
	MinutesToFullCharge         int         `json:"minutes_to_full_charge"`
	NotEnoughPowerToHeat        bool        `json:"not_enough_power_to_heat"`
	ScheduledChargingPending    bool        `json:"scheduled_charging_pending"`
	ScheduledChargingStartTime  int         `json:"scheduled_charging_start_time"`
	TimeToFullCharge            float64     `json:"time_to_full_charge"`
	Timestamp                   int64       `json:"timestamp"`
	TripCharging                bool        `json:"trip_charging"`
	UsableBatteryLevel          int         `json:"usable_battery_level"`
	UserChargeEnableRequest     interface{} `json:"user_charge_enable_request"`
}

func (c ChargeState) GetState() string {
	return c.ChargingState
}

func (c ChargeState) GetMaxAmps() int {
	return c.ChargeCurrentRequestMax
}

type ChargeSiteRes struct {
	BaseRes
	Response struct {
		CongestionSyncTimeUtcSecs int          `json:"congestion_sync_time_utc_secs"`
		DestinationCharging       []ChargeSite `json:"destination_charging"`
		Superchargers             []ChargeSite `json:"superchargers"`
	} `json:"response"`
}

type ChargeSite struct {
	Location struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	} `json:"location"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	DistanceMiles   float64 `json:"distance_miles"`
	AvailableStalls int     `json:"available_stalls"`
	TotalStalls     int     `json:"total_stalls"`
	SiteClosed      bool    `json:"site_closed"`
}

// ChargeState get vehicle charge state
func (t *TeslaApi) ChargeState() (cd *ChargeState, err error) {
	cd = &ChargeState{}
	r := struct {
		BaseRes
		Response ChargeState `json:"response"`
	}{}
	err = t.sendDataRequest("charge_state", &r)
	if err != nil {
		return cd, err
	}
	cd = &r.Response
	t.activeVehicleData.ChargeState = *cd
	return cd, err
}

// NearByChargingSites return list of superchargers and destination chargers
func (t TeslaApi) NearByChargingSites() (cs []ChargeSite, err error) {
	if t.activeVehicle.Id == 0 {
		return cs, ErrNoActiveVehicle
	}
	u := joinPath(commandUrlBase, vehicleEndpoint, t.activeVehicle.GetIdStr(), "nearby_charging_sites")
	res, err := t.apiRequest(http.MethodGet, u, nil)
	if err != nil {
		return cs, err
	}
	r := ChargeSiteRes{}
	err = parseResp(res, &r)
	if err != nil {
		return cs, errors.New(r.Err)
	}
	cs = append(cs, r.Response.DestinationCharging...)
	cs = append(cs, r.Response.Superchargers...)
	return cs, err
}

func (t TeslaApi) ChargeDoorOpen() (o *CommandsRes, err error) {
	return t.sendCommand(cmdChargeDoorOpen, "")
}

func (t TeslaApi) ChargeDoorClose() (o *CommandsRes, err error) {
	return t.sendCommand(cmdChargeDoorClose, "")
}

func (t TeslaApi) ChargeStart() (o *CommandsRes, err error) {
	return t.sendCommand(cmdChargeStart, "")
}

func (t TeslaApi) ChargeStop() (o *CommandsRes, err error) {
	return t.sendCommand(cmdChargeStop, "")
}

func (t TeslaApi) ChargeStandard() (o *CommandsRes, err error) {
	return t.sendCommand(cmdChargeStandard, "")
}

func (t TeslaApi) ChargeMaxRange() (o *CommandsRes, err error) {
	return t.sendCommand(cmdChargeMax, "")
}

func (t TeslaApi) SetChargeLimit(percent int) (o *CommandsRes, err error) {
	if percent > 100 || percent < 0 {
		return o, ErrInvalidChargePercent
	}
	return t.sendCommand(
		cmdChargeLimit,
		t.formUrlEncode(map[string]string{
			"percent": strconv.FormatInt(int64(percent), 10),
		}),
	)
}

func (t TeslaApi) SetChargeAmps(amps int) (o *CommandsRes, err error) {
	if amps < 1 || amps > 1000 ||
		(t.activeVehicleData.ChargeState.GetMaxAmps() != 0 && amps > t.activeVehicleData.ChargeState.GetMaxAmps()) {
		return o, ErrInvalidAmps
	}
	return t.sendCommand(
		cmdChargeAmp,
		t.formUrlEncode(map[string]string{
			"charging_amps": strconv.FormatInt(int64(amps), 10),
		}),
	)
}

func (t TeslaApi) SetScheduledCharge(enable bool, time int) (o *CommandsRes, err error) {
	if time < 0 || time > 1440 {
		return o, ErrInvalidChargeTime
	}
	return t.sendCommand(
		cmdChargeSchedule,
		t.formUrlEncode(map[string]string{
			"enable": strconv.FormatBool(enable),
			"time":   strconv.FormatInt(int64(time), 10),
		}),
	)
}

// SetScheduledDeparture
// departureTime in min
// end_off_peak_time in min
func (t TeslaApi) SetScheduledDeparture(enable bool, departureTime, endOffPeakTime int, preconditioningEnabled,
	preconditioningWeekdaysOnly, offPeakChargingEnabled, offPeakChargingWeekdaysOnly bool) (o *CommandsRes, err error) {
	if departureTime < 0 || departureTime > 1440 {
		return o, ErrInvalidChargeTime
	}
	if endOffPeakTime < 0 || endOffPeakTime > 1440 {
		return o, ErrInvalidChargeTime
	}
	return t.sendCommand(
		cmdChargeDeparture,
		t.formUrlEncode(map[string]string{
			"enable":                          strconv.FormatBool(enable),
			"departure_time":                  strconv.FormatInt(int64(departureTime), 10),
			"preconditioning_enabled":         strconv.FormatBool(preconditioningEnabled),
			"preconditioning_weekdays_only":   strconv.FormatBool(preconditioningWeekdaysOnly),
			"off_peak_charging_enabled":       strconv.FormatBool(offPeakChargingEnabled),
			"off_peak_charging_weekdays_only": strconv.FormatBool(offPeakChargingWeekdaysOnly),
			"end_off_peak_time":               strconv.FormatInt(int64(endOffPeakTime), 10),
		}),
	)
}

func (t TeslaApi) IsCharging() bool {
	return strings.ToLower(t.activeVehicleData.ChargeState.GetState()) == "charging"
}

type SuperChargingHistory struct {
	Vin                 string      `json:"vin"`
	ChargeSessionId     string      `json:"chargeSessionId"`
	SiteLocationName    string      `json:"siteLocationName"`
	ChargeStartDateTime time.Time   `json:"chargeStartDateTime"`
	ChargeStopDateTime  time.Time   `json:"chargeStopDateTime"`
	UnlatchDateTime     time.Time   `json:"unlatchDateTime"`
	CountryCode         string      `json:"countryCode"`
	Credit              interface{} `json:"credit"`
	DisputeDetails      interface{} `json:"disputeDetails"`
	Fees                []struct {
		SessionFeeId  int     `json:"sessionFeeId"`
		FeeType       string  `json:"feeType"`
		CurrencyCode  string  `json:"currencyCode"`
		PricingType   string  `json:"pricingType"`
		RateBase      float64 `json:"rateBase"`
		RateTier1     int     `json:"rateTier1"`
		RateTier2     int     `json:"rateTier2"`
		RateTier3     int     `json:"rateTier3"`
		RateTier4     int     `json:"rateTier4"`
		UsageBase     int     `json:"usageBase"`
		UsageTier1    int     `json:"usageTier1"`
		UsageTier2    int     `json:"usageTier2"`
		UsageTier3    int     `json:"usageTier3"`
		UsageTier4    int     `json:"usageTier4"`
		TotalBase     float64 `json:"totalBase"`
		TotalTier1    int     `json:"totalTier1"`
		TotalTier2    int     `json:"totalTier2"`
		TotalTier3    int     `json:"totalTier3"`
		TotalTier4    int     `json:"totalTier4"`
		TotalDue      float64 `json:"totalDue"`
		NetDue        float64 `json:"netDue"`
		Uom           string  `json:"uom"`
		IsPaid        bool    `json:"isPaid"`
		Status        string  `json:"status"`
		ProcessFlagId int     `json:"processFlagId"`
	} `json:"fees"`
	BillingType string `json:"billingType"`
	Invoices    []struct {
		FileName    string `json:"fileName"`
		ContentId   string `json:"contentId"`
		InvoiceType string `json:"invoiceType"`
		BeInvoiceId string `json:"beInvoiceId"`
		ProcessFlag int    `json:"processFlag"`
	} `json:"invoices"`
	FapiaoDetails   interface{} `json:"fapiaoDetails"`
	ProgramType     string      `json:"programType"`
	VehicleMakeType string      `json:"vehicleMakeType"`
}

func (t TeslaApi) GetSuperChargingHistory() (ch []SuperChargingHistory, err error) {
	res, err := t.teslaAcctApi(http.MethodGet, "charging/api/history", t.formUrlEncode(map[string]string{
		"vin": t.activeVehicleData.Vin,
	}))
	if err != nil {
		return ch, err
	}
	for _, v := range res.Data.([]interface{}) {
		h, ok := v.(SuperChargingHistory)
		if !ok {
			continue
		}
		ch = append(ch, h)
	}
	return ch, err
}

type ChargeHistory struct {
	ScreenTitle  string `json:"screen_title"`
	TotalCharged struct {
		Value          string `json:"value"`
		RawValue       int    `json:"raw_value"`
		AfterAdornment string `json:"after_adornment"`
		Title          string `json:"title"`
	} `json:"total_charged"`
	TotalChargedBreakdown struct {
		Home struct {
			Value          string `json:"value"`
			AfterAdornment string `json:"after_adornment"`
			SubTitle       string `json:"sub_title"`
		} `json:"home"`
		SuperCharger struct {
			Value          string `json:"value"`
			AfterAdornment string `json:"after_adornment"`
			SubTitle       string `json:"sub_title"`
		} `json:"super_charger"`
		Other struct {
			Value          string `json:"value"`
			RawValue       int    `json:"raw_value"`
			AfterAdornment string `json:"after_adornment"`
			SubTitle       string `json:"sub_title"`
		} `json:"other"`
	} `json:"total_charged_breakdown"`
	ChargingHistoryGraph ChargingHistoryGraph `json:"charging_history_graph"`
}

type ChargingHistoryGraph struct {
	DataPoints []struct {
		Timestamp struct {
			Timestamp struct {
				Seconds int `json:"seconds"`
			} `json:"timestamp"`
			DisplayString string `json:"display_string"`
		} `json:"timestamp"`
		Values []struct {
			Value          string `json:"value"`
			RawValue       int    `json:"raw_value,omitempty"`
			AfterAdornment string `json:"after_adornment"`
			SubTitle       string `json:"sub_title,omitempty"`
		} `json:"values"`
	} `json:"data_points"`
	Period struct {
		StartTimestamp struct {
			Seconds int `json:"seconds"`
		} `json:"start_timestamp"`
		EndTimestamp struct {
			Seconds int `json:"seconds"`
		} `json:"end_timestamp"`
	} `json:"period"`
	Interval int `json:"interval"`
	XLabels  []struct {
		Value    string `json:"value"`
		RawValue int    `json:"raw_value"`
	} `json:"x_labels"`
	YLabels []struct {
		Value          string  `json:"value"`
		RawValue       float64 `json:"raw_value,omitempty"`
		AfterAdornment string  `json:"after_adornment,omitempty"`
	} `json:"y_labels"`
	HorizontalGridLines []float64 `json:"horizontal_grid_lines"`
	VerticalGridLines   []float64 `json:"vertical_grid_lines"`
	DiscreteX           bool      `json:"discrete_x"`
	YRangeMax           int       `json:"y_range_max"`
}

func (t TeslaApi) GetChargeHistory() (ch *ChargeHistory, err error) {
	if t.activeVehicle.Id == 0 {
		return ch, ErrNoActiveVehicle
	}
	u := joinPath(commandUrlBase, vehicleEndpoint, t.activeVehicle.GetIdStr(), "charge_history")
	res, err := t.apiRequest(http.MethodGet, u, nil)
	if err != nil {
		return ch, err
	}
	chRes := struct {
		BaseRes
		Response ChargeHistory `json:"response"`
	}{}
	err = parseResp(res, &chRes)
	if err != nil {
		return ch, err
	}
	ch = &chRes.Response
	return ch, err
}
