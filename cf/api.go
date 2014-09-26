package cf

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

type Token struct {
	Access   string `json:"access_token"`
	Refresh  string `json:"refresh_token"`
	Expire   int64  `json:"expires_in"`
	Obtained time.Time
}

type Api struct {
	logEnabled     bool
	logBodyEnabled bool
	domain         string
	username       string
	password       string
	token          Token
}

func NewApi(dom string, user string, pass string, log bool, logBody bool) Api {

	if dom == "" {
		panic("Need a DOMAIN to point the request somewhere...")
	}

	return Api{domain: dom, username: user, password: pass, logEnabled: log, logBodyEnabled: logBody, token: Token{}}
}

func (api *Api) Put(url string, body *bytes.Buffer) (*http.Response, error) {

	u := fmt.Sprintf("https://api.%s%s", api.domain, url)
	// u := fmt.Sprintf("http://localhost:9999/%s/%s", api.domain, url)

	req, err := http.NewRequest("PUT", u, body)
	req.Body.Close()
	if err != nil {
		return nil, err
	}
	return api.doAuthenticatedRequest(req)
}

func (api *Api) Get(url string) (*http.Response, error) {

	u := fmt.Sprintf("https://api.%s%s", api.domain, url)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return api.doAuthenticatedRequest(req)
}

func (api *Api) doAuthenticatedRequest(req *http.Request) (*http.Response, error) {

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.getAccessToken())

	return api.doRequest(req)
}

func (api *Api) doRequest(req *http.Request) (*http.Response, error) {

	if api.logEnabled {
		reqString, _ := httputil.DumpRequest(req, api.logBodyEnabled)
		fmt.Printf("\n\n=====> Req: \n%s", string(reqString))
	}

	res, err := api.newClient().Do(req)
	// defer res.Body.Close()

	if api.logEnabled {
		resString, _ := httputil.DumpResponse(res, api.logBodyEnabled)
		if err != nil {
			fmt.Printf("-----> Resp ERROR: \n%s - %s", string(resString), err.Error())

		} else {
			fmt.Printf("-----> Resp: \n%s\n\n", string(resString))
		}
	}

	return res, err
}

func (api *Api) newClient() *http.Client {

	//Skip self signed cert validation...
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}

func (api *Api) getAccessToken() string {

	//If we have not token, get the first one
	if api.token.Access == "" {

		fmt.Println("[Token] Will try to get a new token")

		url := fmt.Sprintf("https://login.%s/oauth/token?username=%s&password=%s&grant_type=password", api.domain, api.username, api.password)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic("Could not try to get a new token: " + err.Error())
		}

		req.Header.Set("Authorization", "Basic Y2Y6") //Magic token! 'cf:' in base64
		res, err := api.doRequest(req)
		if err != nil {
			panic("could not get a new token: " + err.Error())
		}

		json.NewDecoder(res.Body).Decode(&api.token)
		api.token.Obtained = time.Now()
		return api.token.Access
	}

	since := time.Now().Unix() - api.token.Obtained.Unix()
	hasNotExpired := since < api.token.Expire-5

	if hasNotExpired {
		return api.token.Access
	}

	fmt.Println("[Token] Will refresh token")
	return api.refreshToken()
}

func (api *Api) refreshToken() string {

	url := fmt.Sprintf("https://login.%s/oauth/token?grant_type=refresh_token&refresh_token=%s", api.domain, api.token.Refresh)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Could not try to get a new token: " + err.Error())
	}

	req.Header.Set("Authorization", "Basic Y2Y6") //Magic token! 'cf:' in base64
	res, err := api.doRequest(req)
	if err != nil {
		panic("could not refresh token: " + err.Error())
	}

	json.NewDecoder(res.Body).Decode(&api.token)
	api.token.Obtained = time.Now()
	// fmt.Println("Token refreshed ok!!!!")
	return api.token.Access
}
