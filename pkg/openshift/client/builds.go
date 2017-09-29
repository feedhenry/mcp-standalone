/**

The build package has the same signature as the openshift client library so hopefully will reduce the pain when moving over to the official client.


*/

package client

import (
	"encoding/json"

	"bytes"
	"fmt"
	"net/http"

	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	"github.com/pkg/errors"
)

type BuildConfigs struct {
	ns          string
	host, token string
	restClient  *http.Client
}

const (
	buildconfigURL = "%s/oapi/v1/namespaces/%s/buildconfigs"
)

func (bc *BuildConfigs) Create(config *build.BuildConfig) (*build.BuildConfig, error) {
	payload, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal buildconfig payload")
	}
	u := fmt.Sprintf(buildconfigURL, bc.host, bc.ns)
	req, err := http.NewRequest("POST", u, bytes.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare post request for buildconfig create")
	}
	req.Header.Set("Authorization", "Bearer "+bc.token)
	req.Header.Set("Content-type", "application/json")
	res, err := bc.restClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to create buildconfig ")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return nil, errors.Wrap(err, "unexpected status creating buildconfig")
	}
	return config, nil
}

func (bd *BuildConfigs) Update(config *build.BuildConfig) (*build.BuildConfig, error) {
	return nil, nil
}
