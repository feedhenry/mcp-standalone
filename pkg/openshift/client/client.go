package client

import (
	"net/http"

	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
)

type Client struct {
	ns          string
	host, token string
	restClient  *http.Client
}

func NewClient(ns, host, token string, httpClient *http.Client) *Client {
	return &Client{
		ns:         ns,
		host:       host,
		token:      token,
		restClient: httpClient,
	}
}

// Builds provides a REST client for Builds
func (c *Client) Builds(namespace string) BuildInterface {
	return &BuildConfigs{ns: namespace, host: c.host, token: c.token, restClient: c.restClient}
}

// BuildConfigs provides a REST client for BuildConfigs
func (c *Client) BuildConfigs(namespace string) BuildConfigInterface {
	return &BuildConfigs{ns: namespace, host: c.host, token: c.token, restClient: c.restClient}
}

// Interface exposes methods on OpenShift resources.
type Interface interface {
	BuildsNamespacer
	BuildConfigsNamespacer
}

type BuildsNamespacer interface {
	Builds(namespace string) BuildInterface
}

type BuildConfigsNamespacer interface {
	BuildConfigs(namespace string) BuildConfigInterface
}

type BuildInterface interface{}

type BuildConfigInterface interface {
	Create(config *build.BuildConfig) (*build.BuildConfig, error)
	Update(config *build.BuildConfig) (*build.BuildConfig, error)
}
