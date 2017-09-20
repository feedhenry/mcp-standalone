package openshift

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

//AuthCheckerBuilder for building AuthCheckers
type AuthCheckerBuilder struct {
	Host          string
	Token         string
	SkipCertCheck bool
}

// AuthChecker checks authorizations against resource in namespaces
type AuthChecker struct {
	Host          string
	Token         string
	SkipCertCheck bool
}

// Build an AuthChecker and return it
func (acb *AuthCheckerBuilder) Build() mobile.AuthChecker {
	return &AuthChecker{
		Host:          acb.Host,
		Token:         acb.Token,
		SkipCertCheck: acb.SkipCertCheck,
	}
}

// IgnoreCerts sets the config to ignore future certificate errors
func (acb *AuthCheckerBuilder) IgnoreCerts() mobile.AuthCheckerBuilder {
	return &AuthCheckerBuilder{
		Host:          acb.Host,
		Token:         acb.Token,
		SkipCertCheck: true,
	}
}

// WithToken stores the provided for creating future AuthCheckers
func (acb *AuthCheckerBuilder) WithToken(token string) mobile.AuthCheckerBuilder {
	return &AuthCheckerBuilder{
		Host:          acb.Host,
		SkipCertCheck: acb.SkipCertCheck,
		Token:         token,
	}
}

type authCheckJsonPayload struct {
	Namespace          string `json:"namespace"`
	Verb               string `json:"verb"`
	ResourceAPIGroup   string `json:"resourceAPIGroup"`
	ResourceAPIVersion string `json:"resourceAPIVersion"`
	Resource           string `json:"resource"`
	ResourceName       string `json:"resourceName"`
	Path               string `json:"path"`
	IsNonResourceURL   string `json:"isNonResourceURL"`
}

// Check that the resource in the provided namespace can be written to by the current user
func (ac *AuthChecker) Check(resource, namespace string) (bool, error) {
	u, err := url.Parse(ac.Host)
	if err != nil {
		return false, errors.Wrap(err, "openshift.ac.Check -> failed to parse openshift host when attempting to read user")
	}
	u.Path = path.Join("/oapi/v1/localsubjectaccessreviews")
	payload := authCheckJsonPayload{
		Namespace:          namespace,
		Verb:               "update",
		ResourceAPIGroup:   "apps/v1beta1",
		ResourceAPIVersion: "v1",
		Resource:           "deployments",
	}
	strPayload, err := json.Marshal(payload)
	if err != nil {
		return false, errors.Wrap(err, "openshift.ac.Check -> failed to build payload for check authorization")
	}
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(string(strPayload)))
	if err != nil {
		return false, errors.Wrap(err, "openshift.ac.Check -> failed to build request to check authorization")
	}
	req.Header.Set("authorization", "bearer "+ac.Token)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ac.SkipCertCheck},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "openshift.ac.Check -> failed to make request to check authorization")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return false, &AuthenticatationError{Message: "openshift.ac.Check -> access was denied", StatusCode: resp.StatusCode}
		}

		return false, errors.New(fmt.Sprintf("openshift.ac.Check -> unexpected response code from openshift %v", resp.StatusCode))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "openshift.ac.Check -> failed to read the response body after reading user")
	}
	fmt.Println(string(data))
	return true, nil
}

// NewAuthCheckerBuilder created and returned with the provided namespace and host
func NewAuthCheckerBuilder(host string) mobile.AuthCheckerBuilder {
	return &AuthCheckerBuilder{
		Host:          host,
		SkipCertCheck: false,
	}
}
