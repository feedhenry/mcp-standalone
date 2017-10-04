package testclient

import (
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"k8s.io/client-go/testing"
)

type Client struct {
	ns          string
	host, token string
	Fake        *testing.Fake
}

func NewClient(ns, host, token string, fake *testing.Fake) *Client {
	return &Client{
		ns:    ns,
		host:  host,
		token: token,
		Fake:  fake,
	}
}

// Builds provides a REST client for Builds
func (c *Client) Builds(namespace string) client.BuildInterface {
	return NewFakeBuilds(namespace, c.Fake)
}

// BuildConfigs provides a REST client for BuildConfigs
func (c *Client) BuildConfigs(namespace string) client.BuildConfigInterface {
	return NewFakeBuildConfigs(namespace, c.Fake)
}
