package k8s

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"encoding/json"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

type Metadata struct {
	Issuer                        string   `json:"issuer"`
	AuthorizationEndpoint         string   `json:"authorization_endpoint"`
	TokenEndpoint                 string   `json:"token_endpoint"`
	ScopesSupported               []string `json:"scopes_supported"`
	ResponseTypesSupported        []string `json:"response_types_supported"`
	GrantTypesSupported           []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
}

// GetK8IssuerHost will parse out the host from the meta response
func (m *Metadata) GetK8IssuerHost() (string, error) {
	parsed, err := url.Parse(m.Issuer)
	if err != nil {
		return "", errors.Wrap(err, "GetK8IssuerHost failed to parse k8s url")
	}
	return parsed.Host, nil
}

func GetMetadata(k8shost string, requester mobile.ExternalHTTPRequester) (*Metadata, error) {
	url := fmt.Sprintf("%s%s", k8shost, "/.well-known/oauth-authorization-server")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request to retrieve OpenShift server metadata")
	}
	res, err := requester.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request to retrieve OpenShift server metadata")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "attempted to read body after failed OpenShift server metadata call")
		}
		if nil != data {
			return nil, errors.Wrap(err, "unexpected response from OpenShift"+string(data))
		}
		return nil, errors.New("unexpected response from OpenShift: " + res.Status)
	}

	metadata := &Metadata{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(metadata); err != nil {
		return nil, errors.Wrap(err, "OpenShift server metadata: Attempted to decode payload")
	}

	if !checkScopes(*metadata) {
		return nil, errors.New(fmt.Sprintf("OpenShift server metadata is missing required scopes. Got %v", metadata.ScopesSupported))
	}

	if !checkResponseTypes(*metadata) {
		return nil, errors.New(fmt.Sprintf("OpenShift server metadata is missing required response types. Got %v", metadata.ResponseTypesSupported))
	}

	if !checkGrantTypes(*metadata) {
		return nil, errors.New(fmt.Sprintf("OpenShift server metadata is missing required grant types. Got %v", metadata.GrantTypesSupported))
	}

	return metadata, nil
}

func checkScopes(metadata Metadata) bool {
	hasUserInfoScope := false
	hasUserCheckAccessScope := false
	for _, b := range metadata.ScopesSupported {
		if !hasUserInfoScope && b == "user:info" {
			hasUserInfoScope = true
		}
		if !hasUserCheckAccessScope && b == "user:check-access" {
			hasUserCheckAccessScope = true
		}
	}
	return hasUserInfoScope && hasUserCheckAccessScope
}

func checkResponseTypes(metadata Metadata) bool {
	hasCodeType := false
	hasTokenType := false
	for _, b := range metadata.ResponseTypesSupported {
		if !hasCodeType && b == "code" {
			hasCodeType = true
		}
		if !hasTokenType && b == "token" {
			hasTokenType = true
		}
	}
	return hasCodeType && hasTokenType
}

func checkGrantTypes(metadata Metadata) bool {
	hasAuthorizationCodeType := false
	for _, b := range metadata.GrantTypesSupported {
		if !hasAuthorizationCodeType && b == "authorization_code" {
			hasAuthorizationCodeType = true
		}
	}
	return hasAuthorizationCodeType
}
