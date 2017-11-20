package jenkins

import (
	"io"

	"net/url"

	"net/http"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

type Client struct {
	externalRequester mobile.ExternalHTTPRequester
	logger            *logrus.Logger
}

func NewClient(externalRequester mobile.ExternalHTTPRequester, logger *logrus.Logger) *Client {
	return &Client{
		externalRequester: externalRequester,
		logger:            logger,
	}
}

type buildArtifact struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

func (ba *buildArtifact) artifactUrl(host string) string {
	return host + ba.URL
}

// Retrieve will make a request to the artifact url and if successful return the body of that request as a readerCloser
func (c *Client) Retrieve(location *url.URL, token string) (io.ReadCloser, error) {
	u := location.String()
	c.logger.Debug("calling jenkins url ", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jenkins download request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := c.externalRequester.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error doing jenkins artifact list request")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected response code from Jenkins download " + res.Status)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	decoder := json.NewDecoder(res.Body)
	buildArtifacts := []*buildArtifact{}
	if err := decoder.Decode(&buildArtifacts); err != nil {
		return nil, errors.Wrap(err, "failed to decode artifact list response from jenkins")
	}
	if len(buildArtifacts) == 0 {
		return nil, errors.New("no artifacts returned for build")
	}
	retURL := buildArtifacts[0].artifactUrl(location.Scheme + "://" + location.Host)
	c.logger.Debug("calling jenkins url for artifact ", retURL)
	req, err = http.NewRequest("GET", retURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jenkins artifact download request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	artResponse, err := c.externalRequester.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error doing jenkins artifact download request")
	}
	if artResponse.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected response code from Jenkins artifact download " + res.Status)
	}
	//not closing the body here as we are handing it off to be streamed and closed
	return artResponse.Body, nil
}
