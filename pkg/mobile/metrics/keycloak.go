package metrics

import (
	"encoding/json"
	"strings"

	"time"

	"fmt"

	"net/http"
	"net/url"

	"io/ioutil"

	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

type Keycloak struct {
	requestBuilder     mobile.HTTPRequesterBuilder
	serviceRepoBuilder mobile.ServiceRepoBuilder
	ServiceName        string
	logger             *logrus.Logger
}

type token struct {
	val        string
	validUntil int64
}

var tokenCache = map[string]*token{}

const timeFomrat = "2006-01-02 15:04:05"

func NewKeycloak(rbuilder mobile.HTTPRequesterBuilder, serviceRepoBuilder mobile.ServiceRepoBuilder, l *logrus.Logger) *Keycloak {
	return &Keycloak{requestBuilder: rbuilder, serviceRepoBuilder: serviceRepoBuilder, logger: l, ServiceName: "keycloak"}
}

// Gather will retrieve varous metrics from keycloak
func (kc *Keycloak) Gather() ([]*metric, error) {
	kc.logger.Debug("keycloak metrics gathering ")
	svc, err := kc.serviceRepoBuilder.UseDefaultSAToken().Build()
	if err != nil {
		return nil, errors.Wrap(err, "keycloak gather failed to create svcruder using default service account token")
	}
	kcServices, err := svc.List(func(attrs mobile.Attributer) bool {
		return attrs.GetName() == kc.ServiceName
	})
	if err != nil {
		return nil, errors.Wrap(err, "keycloak gather failed to list existing services")
	}
	if len(kcServices) == 0 {
		return nil, &noServiceProvisionedErr{Message: " no keycloak service present in namespace "}
	}
	kcService := kcServices[0] //TODO deal with more than one
	host := kcService.Host
	username := kcService.Params["admin_username"]
	pass := kcService.Params["admin_password"]
	realm := kcService.Params["realm"]
	token, err := kc.getToken(host, username, pass, realm)
	if err != nil {
		return nil, err
	}
	cs, err := kc.getClientStats(host, token, realm)
	if err != nil {
		kc.logger.Error("keycloak: failed to get client stats ", err)
	}
	events, err := kc.getRealmEvents(host, token, realm)
	if err != nil {
		kc.logger.Error("keycloak: failed to get realm events ", err)
	}
	var kcMetrics = []*metric{}
	if len(cs) > 0 {
		clientMetrics := processClientStats(cs)
		kcMetrics = append(kcMetrics, clientMetrics...)
	}
	if len(events) > 0 {
		eventMetrics := processRealmEvents(events)
		kcMetrics = append(kcMetrics, eventMetrics...)
	}
	return kcMetrics, nil
}

func processClientStats(stats []*clientStat) []*metric {
	now := time.Now().Format(timeFomrat)
	ret := []*metric{}
	for _, s := range stats {
		active, _ := strconv.ParseInt(s.Active, 10, 0)
		m := &metric{
			Type:   s.ClientID,
			XValue: now,
			YValue: active,
		}
		ret = append(ret, m)
	}
	return ret
}

func processRealmEvents(events []*eventType) []*metric {
	ret := []*metric{}
	now := time.Now().Format(timeFomrat)
	for _, e := range events {
		added := false
		for i := range ret {
			existing := ret[i]
			if existing.Type == e.Type {
				existing.YValue++
				added = true
				break
			}
		}
		if !added {
			ret = append(ret, &metric{Type: e.Type, XValue: now, YValue: 1})
		}
	}
	return ret
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	if resp.StatusCode != 200 {
		return "", errors.New("failed to login to keycloak response code was: " + resp.Status + " url called was : " + u)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "keylcloak: failed to read response ")
	}
	payload := map[string]interface{}{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", errors.Wrap(err, "failed to Unmarshal keycloak response ")
	}

	validFor, ok := payload["expires_in"].(float64)
	if !ok {
		return "", errors.New("payload expires_in failed to convert to float64 or was not present ")
	}
	accessToken, ok := payload["access_token"].(string)
	if !ok {
		return "", errors.New("payload access_token failed to convert to string or was not present ")
	}
	tokenCache[realm] = &token{val: accessToken, validUntil: tokenRequestTime + int64(validFor) - 2} //give a bit of a margin of error
	return accessToken, nil
}

//{"clientId":"account","active":"1","id":"fad0b64e-818e-4545-8b25-6a32e80c8484"
type clientStat struct {
	ClientID string `json:"clientID"`
	Active   string `json:"active"`
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
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	decode := json.NewDecoder(res.Body)
	clientStats := []*clientStat{}
	if err := decode.Decode(&clientStats); err != nil {
		return nil, errors.Wrap(err, "failed to decode client stats from keycloak")
	}
	return clientStats, nil

}

type eventType struct {
	Type string `json:"type"`
}

func (kc *Keycloak) getRealmEvents(host, token, realm string) ([]*eventType, error) {
	//url /admin/realms/{realm}/admin-events
	u := fmt.Sprintf("%s/auth/admin/realms/%s/events", host, realm)
	httpClient := kc.requestBuilder.Insecure(true).Timeout(10).Build()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to crate get request for events ")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request to get realm events")
	}
	if res.StatusCode != 200 {
		return nil, errors.New("unexpected response code when getting realm events expected 200 but got: " + res.Status)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	var events = []*eventType{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&events); err != nil {
		return nil, errors.Wrap(err, "keycloak: failed to decode events response")
	}
	return events, nil
}
