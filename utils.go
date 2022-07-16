package tesla

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

const (
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	commandUrlBase  = "https://owner-api.teslamotors.com"
	teslaUrlBase    = "https://www.tesla.com"
	accountEndpoint = "teslaaccount"
	vehicleEndpoint = "/api/1/vehicles"
	userEndpoint    = "/api/1/users"
	dataReqEndpoint = "data_request"
	commandEndpoint = "command"
)

var (
	ErrVehicleInSleep  = errors.New(`vehicle unavailable: {:error=>"vehicle unavailable:"}`)
	ErrInvalidToken    = errors.New("invalid bearer token")
	ErrNoActiveVehicle = errors.New("no active vehicle please run SetActiveVehicle")
	ErrWakeTimeout     = errors.New("wake up vehicle timed out")
	WakeTimeoutSec     = 30
	// AutoWakeUp Automatically wake up vehicle if vehicle in sleep
	AutoWakeUp               = true
	defaultReqInterval       = time.Second * 5
	ChargeStateReqInterval   = defaultReqInterval
	DriveStateReqInterval    = defaultReqInterval
	ClimateStateReqInterval  = defaultReqInterval
	VehicleStateReqInterval  = defaultReqInterval
	VehicleConfigReqInterval = defaultReqInterval
	GuiSettingReqInterval    = defaultReqInterval
)

type BaseRes struct {
	Err    string `json:"error,omitempty"`
	ErrDes string `json:"error_description,omitempty"`
	ErrUri string `json:"error_uri,omitempty"`
	// fields for tesla.com/teslaaccount response
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Success bool   `json:"success,omitempty"`
}

type CommandsRes struct {
	BaseRes
	Response struct {
		Reason string `json:"reason"`
		Result bool   `json:"result"`
	} `json:"response"`
}

type TeslaAcctRes struct {
	BaseRes
	Data interface{} `json:"data"`
}

type JwtAccessClaims struct {
	Aud  interface{}            `json:"aud"`
	Azp  string                 `json:"azp"`
	Sub  string                 `json:"sub"`
	Scp  []string               `json:"scp"`
	Amr  []string               `json:"amr"`
	Data map[string]interface{} `json:"data,omitempty"`
	jwt.StandardClaims
}

func (t *TeslaApi) WaitForWakeUp(timeout int) (err error) {
	if timeout <= 0 {
		timeout = WakeTimeoutSec
	}
	if t.activeVehicle.IsInSleep() {
		_, err = t.WakeUp()
		if err != nil {
			return err
		}
	}
	c := make(chan error, 1)
	ready := make(chan error, 1)
	go func() {
		for {
			v, err := t.ListVehicleById(t.activeVehicle.GetIdStr())
			if err != nil {
				ready <- err
				return
			}
			t.activeVehicle = v
			if !v.IsInSleep() {
				ready <- nil
				return
			}
			time.Sleep(time.Second * 5)
		}
	}()

	select {
	case err = <-c:
		return err
	case <-time.After(time.Second * time.Duration(timeout)):
		c <- ErrWakeTimeout
	}
	return err
}

func (t *TeslaApi) apiRequest(method, url string, body io.Reader) (res *http.Response, err error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("User-Agent", "")
	req.Header.Add("X-Tesla-User-Agent", "")
	if method == http.MethodPost {
		if url == tokenUrl {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req.Header.Add("Content-Type", "application/json")
		}
		req.Header.Add("Content-Type", "application/json")
	}
	if t.accessToken != "" {
		expired := isTokenExpired(t.accessToken)
		if expired && t.refreshToken != "" && !t.renewingToken {
			log.Infof("access token expired renewing with refresh token")
			t.renewingToken = true
			err = t.renewToken()
			if err != nil {
				return res, err
			}
		} else if expired && !t.renewingToken {
			return res, errors.New("access token expired but no refresh token provided")
		}
		req.Header.Add("Authorization", "Bearer "+t.accessToken)
		if t.refreshToken != "" {
			t.setAuthCookies()
		}
	}
	for _, k := range t.cookies {
		req.AddCookie(k)
	}
	res, err = t.client.Do(req)
	if err != nil {
		return res, err
	}
	// Wake up vehicle
	if res.StatusCode == 408 && AutoWakeUp {
		err = t.WaitForWakeUp(WakeTimeoutSec)
		if err == nil {
			return t.apiRequest(method, url, body)
		}
	}
	t.setCookies(req.Cookies())
	return res, err
}

func (t *TeslaApi) sendDataRequest(reqDataType string, reqResp interface{}) (err error) {
	if t.activeVehicle.Id == 0 {
		return ErrNoActiveVehicle
	}
	u := joinPath(commandUrlBase, vehicleEndpoint, t.activeVehicle.GetIdStr(), dataReqEndpoint, reqDataType)
	res, err := t.apiRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	err = parseResp(res, reqResp)
	if err != nil {
		return err
	}
	return err
}

func (t *TeslaApi) sendCommand(command, body string) (cmdRes *CommandsRes, err error) {
	cmdRes = &CommandsRes{}
	if t.activeVehicle.Id == 0 {
		return cmdRes, ErrNoActiveVehicle
	}
	u := joinPath(commandUrlBase, vehicleEndpoint, t.activeVehicle.GetIdStr(), commandEndpoint, command)
	res, err := t.apiRequest(http.MethodPost, u, strings.NewReader(body))
	if err != nil {
		return cmdRes, err
	}
	c := CommandsRes{}
	err = parseResp(res, &c)
	if err != nil {
		return &c, errors.New(c.Err)
	}
	cmdRes = &c
	return cmdRes, err
}

func (t *TeslaApi) teslaAcctApi(method, endpoint, body string) (r *TeslaAcctRes, err error) {
	u := joinPath(teslaUrlBase, accountEndpoint, endpoint)
	res, err := t.apiRequest(method, u, strings.NewReader(body))
	if err != nil {
		return r, err
	}
	r = &TeslaAcctRes{}
	err = parseResp(res, r)
	if err != nil {
		return r, errors.New(r.Message)
	}
	return r, err
}

func (t TeslaApi) formUrlEncode(form map[string]string) (content string) {
	f := url.Values{}
	for k, v := range form {
		f.Add(k, v)
	}
	content = f.Encode()
	return content
}

func (t TeslaApi) jsonEncode(data interface{}) (content string) {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}

func getAuthBase(global bool) string {
	if global {
		return authGlobalBase
	}
	return authCnBase
}

func randomStr(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func UnmarshallJWT(token string) (result *JwtAccessClaims) {
	t, _, err := new(jwt.Parser).ParseUnverified(token, &JwtAccessClaims{})
	if err != nil {
		log.Errorf("failed to parse jwt token: %s", err)
	}
	if claims, ok := t.Claims.(*JwtAccessClaims); ok {
		return claims
	}
	return result
}

func isTokenExpired(token string) (expired bool) {
	expired = true
	jwtClaims := UnmarshallJWT(token)
	if jwtClaims == nil {
		return expired
	}
	exp := jwtClaims.ExpiresAt
	if time.Since(time.Unix(exp, 0)) < 0 {
		return false
	}
	return expired
}

func joinPath(baseUrl string, paths ...string) string {
	u, _ := url.Parse(baseUrl)
	for _, p := range paths {
		u.Path = path.Join(u.Path, p)
	}
	return u.String()
}

func parseResp(res *http.Response, respType interface{}) (err error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decoder.UseNumber()
	err = decoder.Decode(respType)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code is %d", res.StatusCode))
	}
	return err
}

// if structKeyMap is not empty will use structKeyMap to find the field name
// outType must be a pointer reference
func convertMapResp(in map[string]interface{}, outType interface{}, structKeyMap map[string]string) {
	r := reflect.ValueOf(outType)
	if r.Kind() != reflect.Ptr {
		log.Errorf("outType is %v not pointer refernce", r.Kind())
		return
	}
	e := r.Elem()
	for key, val := range in {
		if val == nil {
			continue
		}
		var f reflect.Value
		if len(structKeyMap) > 0 {
			structField, ok := structKeyMap[key]
			if !ok {
				continue
			}
			f = e.FieldByName(structField)
		} else {
			f = e.FieldByName(key)
		}

		if f.IsValid() && f.CanSet() {
			kind := f.Kind()
			switch kind {
			case reflect.Float64:
				if num, ok := val.(float64); ok {
					f.SetFloat(num)
				} else {
					n, err := val.(json.Number).Float64()
					if err != nil {
						log.Warnf("failed to convert val %v to float64", val)
					}
					f.SetFloat(n)
				}
			case reflect.Int64, reflect.Int:
				if num, ok := val.(int); ok {
					f.SetInt(int64(num))
				} else {
					n, err := val.(json.Number).Int64()
					if err != nil {
						log.Warnf("failed to convert val %v to int", val)
					}
					f.SetInt(n)
				}
			case reflect.String:
				f.SetString(val.(string))
			case reflect.Slice:
				if reflect.TypeOf(val) == reflect.TypeOf([]string{}) {
					for _, v := range val.([]string) {
						f.SetString(v)
					}
				}
			case reflect.Bool:
				f.SetBool(val.(bool))
			}
		}
	}
}

func timestampSince(timestamp int64) time.Duration {
	return time.Since(time.UnixMilli(timestamp))
}
