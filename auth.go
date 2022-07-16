package tesla

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	log "github.com/sirupsen/logrus"
)

const (
	scopeField               = "scope"
	grantTypeField           = "grant_type"
	clientIdField            = "client_id"
	clientSecretField        = "client_secret"
	refreshTokenField        = "refresh_token"
	codeVerifierField        = "code_verifier"
	codeField                = "code"
	stateField               = "state"
	responseTypeField        = "response_type"
	redirectUriField         = "redirect_uri"
	codeChallengeField       = "code_challenge"
	codeChallengeMethodField = "code_challenge_method"
	grantTypeRefresh         = "refresh_token"
	grantTypeAuthCode        = "authorization_code"
	clientSecret             = "c7257eb71a564034f9419ee651c7d0e5f7aa6bfbd18bafb5c5c033b093bb2fa3"
	webAuthCookieField       = "authTeslaWebToken"
	webRefreshCookieField    = "teslaSSORefreshToken"
)

var (
	ownerClientId       = "ownerapi"
	authGlobalBase      = "https://auth.tesla.com"
	authCnBase          = "https://auth.tesla.cn"
	authEndpoint        = "/oauth2/v3/authorize"
	tokenEndpoint       = "/oauth2/v3/token"
	authUrl             = authGlobalBase + authEndpoint
	tokenUrl            = authGlobalBase + tokenEndpoint
	redirectUri         = "https://auth.tesla.com/void/callback"
	codeChallengeMethod = "S256"
	scope               = "openid email offline_access"
	accessTokenExpire   = time.Now().Add(time.Hour * 8)
	refreshTokenExpire  = time.Now().Add(time.Hour * 8760)
)

type AuthRes struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	State        string `json:"state,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	CreatedAt    int    `json:"created_at,omitempty"`
	IdToken      string `json:"id_token,omitempty"`
	BaseRes
}

type AuthReq struct {
	AuthBaseUrl   string
	CodeVerifier  string
	CodeChallenge string
	identity      string
	credential    string
}

// if customVerifier is empty generate random verifier
func (t TeslaApi) getVerifier(customVerifier string) string {
	if customVerifier != "" && len(customVerifier) == 86 {
		return customVerifier
	} else if len(customVerifier) != 86 {
		log.Warnf("%s is not 86 length generating new randome verifier", customVerifier)
	}
	// fix 86 length for verifier
	b := randomStr(86)
	return string(b)
}

func (t TeslaApi) getChallenge() string {
	return base64.StdEncoding.EncodeToString([]byte(t.CodeVerifier))
}

// Step 1: Obtain the login page
// Step 2: Obtain an authorization code
func (t *TeslaApi) getAuthCode() (authCode string) {
	authCode = ""
	if t.AuthReq.identity == "" && t.AuthReq.credential == "" {
		return authCode
	}
	req, _ := http.NewRequest("", authUrl, nil)
	q := req.URL.Query()
	q.Add(clientIdField, ownerClientId)
	q.Add(codeChallengeField, t.CodeChallenge)
	q.Add(codeChallengeMethodField, codeChallengeMethod)
	q.Add(redirectUriField, redirectUri)
	q.Add(responseTypeField, "code")
	q.Add(scopeField, scope)
	q.Add(stateField, string(randomStr(12)))
	req.URL.RawQuery = q.Encode()

	// open tesla log in page in browser
	l, _ := launcher.New().Headless(false).Launch()
	browser := rod.New()
	if browser != nil {
		defer browser.Close()
	}
	page := browser.
		ControlURL(l).
		SlowMotion(1 * time.Second).
		MustConnect().
		MustPage(authUrl + "?" + req.URL.RawQuery)
	time.Sleep(3 * time.Second)
	page.MustWaitLoad()
	log.Infof("waiting for login page to load")
	page.Sleeper(rod.DefaultSleeper).MustElement("#form-input-identity").MustInput(t.identity)
	page.MustElement("#form-submit-continue").MustClick()

	// Re-enter login
	page.MustWaitLoad()
	log.Infof("waiting for email input to load")
	page.MustElement("#form-input-identity").MustInput(t.identity)
	page.MustElement("#form-submit-continue").MustClick()

	page.MustWaitLoad()
	log.Infof("waiting for credential input to load")
	page.MustElement("#form-input-credential").MustInput(t.credential)
	page.MustElement("#form-submit-continue").MustClick()
	page.MustWaitLoad()

	log.Infof("checking recaptcha")
	if _, err := page.Sleeper(rod.NotFoundSleeper).Element("div.recaptcha.tds-form-item.tds-form-item--error"); err == nil {
		page.MustElement("#form-input-credential").MustInput(t.credential)
		log.Infof("capcha found please complete the capcha in browser and sign in")
		err = waitForSignIn(page)
	}
	i, _ := page.Info()
	if !strings.HasPrefix(i.URL, redirectUri) {
		log.Errorf("not redirectoring to call back page")
		return ""
	}

	r, _ := http.NewRequest("", i.URL, nil)
	authCode = r.URL.Query().Get("code")

	cks, _ := page.Browser().GetCookies()
	var httpCks []*http.Cookie

	// Convert browser cookies to http cookies
	for _, ck := range cks {
		ss := http.SameSiteDefaultMode
		switch ck.SameSite {
		case proto.NetworkCookieSameSiteStrict:
			ss = http.SameSiteStrictMode
		case proto.NetworkCookieSameSiteLax:
			ss = http.SameSiteLaxMode
		default:
			ss = http.SameSiteDefaultMode
		}
		httpCks = append(httpCks, &http.Cookie{
			Name:     ck.Name,
			Value:    ck.Name,
			Path:     ck.Path,
			Domain:   ck.Domain,
			Expires:  ck.Expires.Time(),
			Secure:   ck.Secure,
			HttpOnly: ck.HTTPOnly,
			SameSite: ss,
		})
	}
	t.cookies = httpCks
	return authCode
}

// Step 3: Exchange authorization code for bearer token
func (t *TeslaApi) getToken(authCode string) (err error) {
	if authCode == "" {
		return errors.New("auth code is empty")
	}
	res, err := t.apiRequest(
		http.MethodPost,
		tokenUrl,
		strings.NewReader(t.formUrlEncode(map[string]string{
			grantTypeField:    grantTypeAuthCode,
			clientIdField:     ownerClientId,
			codeField:         authCode,
			codeVerifierField: t.CodeVerifier,
			redirectUriField:  redirectUri,
		})),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	aRes, err := parseAuthRes(res.Body)
	if res.StatusCode != 200 {
		return errors.New(aRes.Err)
	}
	t.setToken(aRes.RefreshToken, aRes.AccessToken)
	t.setAuthCookies()
	t.renewingToken = false
	return err
}

// Step 4: Exchange bearer token for access token
func (t *TeslaApi) renewToken() (err error) {
	if t.refreshToken == "" {
		return errors.New("refresh token is empty")
	}
	res, err := t.apiRequest(
		http.MethodPost,
		tokenUrl,
		strings.NewReader(t.formUrlEncode(map[string]string{
			grantTypeField:    grantTypeRefresh,
			scopeField:        scope,
			clientIdField:     ownerClientId,
			clientSecretField: clientSecret,
			refreshTokenField: t.refreshToken,
		},
		)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	aRes, err := parseAuthRes(res.Body)
	if res.StatusCode != 200 {
		return errors.New(aRes.Err)
	}

	t.setToken(aRes.RefreshToken, aRes.AccessToken)
	t.setAuthCookies()
	t.renewingToken = false
	return err
}

// RefreshToken Return stored refresh token
func (t TeslaApi) RefreshToken() string {
	return t.refreshToken
}

func (t TeslaApi) isAuth() bool {
	return t.accessToken != "" && !isTokenExpired(t.accessToken)
}

func (t *TeslaApi) setToken(refreshToken, accessToken string) {
	t.refreshToken = refreshToken
	t.accessToken = accessToken
	accessTokenExpire = time.Now().Add(time.Hour * 8)
}

func (t *TeslaApi) setCookies(cks []*http.Cookie) {
	cksMap := map[string]*http.Cookie{}
	var newCks []*http.Cookie

	for _, ck := range cks {
		cksMap[ck.Name] = ck
	}

	for _, ock := range t.cookies {
		if time.Since(ock.Expires) < 0 {
			continue
		} else if nck, ok := cksMap[ock.Name]; ok && nck.Domain == ock.Domain {
			newCks = append(newCks, nck)
			delete(cksMap, ock.Name)
		} else {
			newCks = append(newCks, ock)
		}
	}

	for _, v := range cksMap {
		newCks = append(newCks, v)
	}

	t.cookies = newCks
}

func (t *TeslaApi) setAuthCookies() {
	authCook := []*http.Cookie{
		{
			Name:    webAuthCookieField,
			Value:   t.accessToken,
			Domain:  "www.tesla.com",
			Path:    "/",
			Expires: accessTokenExpire,
		},
		{
			Name:    webRefreshCookieField,
			Value:   t.refreshToken,
			Domain:  "www.tesla.com",
			Path:    "/",
			Expires: refreshTokenExpire,
		},
	}
	t.setCookies(authCook)
}

func parseAuthRes(r io.Reader) (res *AuthRes, err error) {
	res = &AuthRes{}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return res, err
	}
	return res, err
}

func waitForSignIn(page *rod.Page) (err error) {
	retry := 30
	start := 0
	for start < retry {
		i, _ := page.Info()
		if strings.HasPrefix(i.URL, redirectUri) {
			return nil
		}
		time.Sleep(time.Second * 2)
		start++
	}
	err = errors.New("captcha complete timed out")
	return err
}
