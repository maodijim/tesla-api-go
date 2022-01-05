package tesla_api_go

import (
	"net/http"
	"reflect"
	"strconv"
)

var (
	vehicleJsonToStructField = map[string]string{
		"id":                        "Id",
		"vehicle_id":                "VehicleId",
		"vin":                       "Vin",
		"display_name":              "DisplayName",
		"option_codes":              "OptionCodes",
		"color":                     "Color",
		"access_type":               "AccessType",
		"tokens":                    "Tokens",
		"state":                     "State",
		"in_service":                "InService",
		"id_s":                      "IdS",
		"calendar_enabled":          "CalendarEnabled",
		"api_version":               "ApiVersion",
		"backseat_token":            "BackseatToken",
		"backseat_token_updated_at": "BackseatTokenUpdatedAt",
	}
)

type VehicleRes struct {
	BaseRes
	Response interface{} `json:"response,omitempty"`
	Count    int         `json:"count,omitempty"`
}

type Vehicle struct {
	Id                     int64    `json:"id"`
	VehicleId              int      `json:"vehicle_id"`
	Vin                    string   `json:"vin"`
	DisplayName            string   `json:"display_name"`
	OptionCodes            string   `json:"option_codes"`
	Color                  string   `json:"color"`
	AccessType             string   `json:"access_type"`
	Tokens                 []string `json:"tokens"`
	State                  string   `json:"state"`
	InService              bool     `json:"in_service"`
	IdS                    string   `json:"id_s"`
	CalendarEnabled        bool     `json:"calendar_enabled"`
	ApiVersion             int      `json:"api_version"`
	BackseatToken          string   `json:"backseat_token"`
	BackseatTokenUpdatedAt int      `json:"backseat_token_updated_at"`
}

func (v Vehicle) IsInSleep() bool {
	return v.State == "asleep" || v.State == ""
}

func (v Vehicle) IsInService() bool {
	return v.InService
}

func (v Vehicle) GetIdStr() string {
	return strconv.FormatInt(v.Id, 10)
}

func (t TeslaApi) ListVehicles() (vs []Vehicle, err error) {
	u := joinPath(commandUrlBase, vehicleEndpoint)
	res, err := t.apiRequest(http.MethodGet, u, nil)
	if err != nil {
		return vs, err
	}
	vRes, err := parseVehicleRes(res)
	if err != nil {
		return vs, err
	}
	vs = vRes.Response.([]Vehicle)
	return vs, err
}

func (t TeslaApi) ListVehicleById(id string) (v Vehicle, err error) {
	u := joinPath(commandUrlBase, vehicleEndpoint, id)
	res, err := t.apiRequest(http.MethodGet, u, nil)
	if err != nil {
		return v, err
	}
	vRes, err := parseVehicleRes(res)
	v = vRes.Response.(Vehicle)
	return v, err
}

func parseVehicleRes(res *http.Response) (vRes *VehicleRes, err error) {
	vRes = &VehicleRes{}
	err = parseResp(res, vRes)
	if err != nil {
		return vRes, err
	}
	if reflect.TypeOf(vRes.Response) == reflect.TypeOf([]interface{}{}) {
		s := vRes.Response.([]interface{})
		var newV []Vehicle
		for _, c := range s {
			newV = append(newV, convertMapToVehicle(c.(map[string]interface{})))
		}
		vRes.Response = newV
	} else {
		vRes.Response = convertMapToVehicle(vRes.Response.(map[string]interface{}))
	}
	return vRes, err
}

func convertMapToVehicle(in map[string]interface{}) (vehicle Vehicle) {
	vehicle = Vehicle{}
	convertMapResp(in, &vehicle, vehicleJsonToStructField)
	return
}
