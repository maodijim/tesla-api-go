package tesla

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/cookiejar"
)

const (
	version = "v0.0.1"
)

type TeslaApi struct {
	AuthReq
	client            *http.Client
	cookies           []*http.Cookie
	activeVehicle     Vehicle
	activeVehicleData VehicleData
	accessToken       string
	refreshToken      string
	renewingToken     bool
}

// Login Log in with username & password will require browser pop up, for No GUI environment use refresh token
func (t TeslaApi) Login() (err error) {
	if t.AuthReq.identity != "" && t.AuthReq.credential != "" {
		err = t.getToken(t.getAuthCode())
	} else if t.refreshToken != "" {
		err = t.renewToken()
	} else {
		return errors.New("username and password or refreshToken is missing")
	}
	return err
}

func (t *TeslaApi) SetActiveVehicle(vehicle Vehicle) (err error) {
	t.activeVehicle = vehicle
	data, err := t.VehicleData()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get vehicle data: %s", err))
	}
	t.activeVehicleData = *data
	return nil
}

// NewTeslaApi when using username and password will log in via browser
func NewTeslaApi(username, password, refreshToken string, global bool) *TeslaApi {
	options := cookiejar.Options{}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	ta := TeslaApi{
		AuthReq: AuthReq{
			AuthBaseUrl: getAuthBase(global),
			identity:    username,
			credential:  password,
		},
		client:       &http.Client{Jar: jar},
		refreshToken: refreshToken,
	}
	authUrl = joinPath(ta.AuthBaseUrl, authEndpoint)
	tokenUrl = joinPath(ta.AuthBaseUrl, tokenEndpoint)
	ta.CodeVerifier = ta.getVerifier("")
	ta.CodeChallenge = ta.getChallenge()
	return &ta
}
