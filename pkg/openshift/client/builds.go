/**

The build package has the same signature as the openshift client library so hopefully will reduce the pain when moving over to the official client.


*/

package client

import (
	"encoding/json"

	"bytes"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuildConfigs struct {
	ns          string
	host, token string
	restClient  *http.Client
}

const (
	buildconfigURL = "%s/oapi/v1/namespaces/%s/buildconfigs"
	getBuildURL    = "%s/oapi/v1/namespaces/%s/builds/%s"
	updateBuildURL = "%s/oapi/v1/namespaces/%s/builds/%s"
	createBuildURL = "%s/oapi/v1/namespaces/%s/buildconfigs/%s/instantiate"
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
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	if res.StatusCode != http.StatusCreated {
		return nil, errors.New("unexpected status creating buildconfig " + res.Status)
	}
	return config, nil
}

type Builds struct {
	ns          string
	host, token string
	restClient  *http.Client
}

func (bc *Builds) Instantiate(name string, buildRequest *build.BuildRequest) (error) {
	u := fmt.Sprintf(createBuildURL, bc.host, bc.ns, name)
	payload, err := json.Marshal(buildRequest)
	if err != nil {
		return errors.Wrap(err, "failed to marshal build request")
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(err, "failed to prepare post request for build")
	}
	req.Header.Set("Authorization", "Bearer "+bc.token)
	req.Header.Set("Content-type", "application/json")
	res, err := bc.restClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make post request for build")
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	if res.StatusCode != http.StatusCreated {
		return errors.New("unexpected status code from creating build " + res.Status)
	}
	build := &build.Build{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(build); err != nil {
		return errors.Wrap(err, "failed to decode build object on create")
	}
	return nil
}

func (bc *Builds) Get(name string, options metav1.GetOptions) (*build.Build, error) {
	u := fmt.Sprintf(getBuildURL, bc.host, bc.ns, name)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare get request for build")
	}
	req.Header.Set("Authorization", "Bearer "+bc.token)
	req.Header.Set("Content-type", "application/json")
	res, err := bc.restClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to get build ")
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code from getting build " + res.Status)
	}
	build := &build.Build{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(build); err != nil {
		return nil, errors.Wrap(err, "failed to decode build object")
	}
	return build, nil
}

func (bc Builds) Update(b *build.Build) (*build.Build, error) {
	u := fmt.Sprintf(updateBuildURL, bc.host, bc.ns, b.Name)
	payload, err := json.Marshal(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal build payload")
	}
	req, err := http.NewRequest("PUT", u, bytes.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare post request for build update")
	}
	req.Header.Set("Authorization", "Bearer "+bc.token)
	req.Header.Set("Content-type", "application/json")
	res, err := bc.restClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to create build ")
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status updating build " + res.Status)
	}
	resbuild := &build.Build{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(resbuild); err != nil {
		return nil, errors.Wrap(err, "failed to decode build object")
	}
	return resbuild, nil
}
