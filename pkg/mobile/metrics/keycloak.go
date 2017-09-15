package metrics

import (
	"encoding/json"
	"math/rand"
	"strings"

	"time"

	"fmt"

	"net/http"
	"net/url"

	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

type Keycloak struct {
	requestBuilder     mobile.HTTPRequesterBuilder
	tokenClientBuilder mobile.TokenScopedClientBuilder
	serviceName        string
	logger             *logrus.Logger
}

type token struct {
	val        string
	validUntil int64
}

var tokenCache = map[string]*token{}

func NewKeycloak(rbuilder mobile.HTTPRequesterBuilder, tokenCBuilder mobile.TokenScopedClientBuilder, serviceName string, l *logrus.Logger) *Keycloak {
	return &Keycloak{requestBuilder: rbuilder, tokenClientBuilder: tokenCBuilder, serviceName: serviceName, logger: l}
}

// Gather will retrieve varous metrics from keycloak
func (kc *Keycloak) Gather() ([]*metric, error) {
	svc, err := kc.tokenClientBuilder.UseDefaultSAToken().MobileServiceCruder("")
	if err != nil {
		return nil, errors.Wrap(err, "keycloack gather failed to create svcruder")
	}
	kcServices, err := svc.List(func(attrs mobile.Attributer) bool {
		return attrs.GetName() == kc.serviceName
	})
	if err != nil {
		return nil, errors.Wrap(err, "keycloack gather failed to list existing services")
	}
	if len(kcServices) == 0 {
		return nil, errors.New(" no keycloak service present")
	}
	kcService := kcServices[0] //TODO deal with more than one
	//TODO get protocol from secret
	host := "http://" + kcService.Host
	username := kcService.Params["admin_username"]
	pass := kcService.Params["admin_password"]
	realm := kcService.Params["realm"]
	now := time.Now()

	token, err := kc.getToken(host, username, pass, realm)
	if err != nil {
		return nil, err
	}

	kc.getClientStats(host, token, realm)

	m := &metric{Type: "logins", XValue: now.Format("2006-01-02 15:04:05"), YValue: rand.Intn(100)}
	return []*metric{m}, nil
}

func (kc *Keycloak) getToken(host, user, pass, realm string) (string, error) {
	if v, ok := tokenCache[realm]; ok && v.validUntil < time.Now().Unix() {
		return v.val, nil
	}
	httpClient := kc.requestBuilder.Insecure(true).Timeout(10).Build()
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", user)
	form.Add("password", pass)
	form.Add("client_id", "admin-cli")
	u := fmt.Sprintf("%s/auth/realms/master/protocol/openid-connect/token", host)
	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "failed to create keycloak request ")
	}
	tokenRequestTime := time.Now().Unix()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to make request to keycloak "+u)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("failed to login to keycloak response code was: " + resp.Status + " url called was : " + u)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response ")
	}
	payload := map[string]interface{}{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", errors.Wrap(err, "failed to Unmarshal listPushApplications ")
	}
	validFor := payload["expires_in"].(float64)
	accessToken := payload["access_token"].(string)
	tokenCache[realm] = &token{val: accessToken, validUntil: tokenRequestTime + int64(validFor) - 2} //give a bit of a margin of error
	return accessToken, nil
}

//{"clientId":"account","active":"1","id":"fad0b64e-818e-4545-8b25-6a32e80c8484"
type clientStat struct {
	clientID string `json:"clientID"`
	active   string `json:"active"`
}

func (kc *Keycloak) getClientStats(host, token, realm string) ([]*clientStat, error) {
	u := fmt.Sprintf("%s/auth/admin/realms/%s/client-session-stats", host, realm)
	httpClient := kc.requestBuilder.Insecure(true).Timeout(10).Build()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request for client-session-stats")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client-session-stats")
	}
	if res.StatusCode != 200 {
		return nil, errors.New("unexpected response code from keycloack server " + res.Status)
	}
	defer res.Body.Close()
	decode := json.NewDecoder(res.Body)
	clientStats := []*clientStat{}
	if err := decode.Decode(&clientStats); err != nil {
		return nil, errors.Wrap(err, "failed to decode client stats from keycloak")
	}
	return clientStats, nil

}

func (kc *Keycloak) getAdminEvents(host, token, realm string) {
	//url /admin/realms/{realm}/admin-events
	u := fmt.Sprintf("%s/auth/admin/realms/%s/events", host, realm)
	fmt.Println("calling ", u)
	httpClient := kc.requestBuilder.Insecure(true).Timeout(10).Build()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {

	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error calling events ", err)
	}
	if res.StatusCode != 200 {
		fmt.Println("unexpected response code fro")
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	fmt.Println(err, string(data))
}
