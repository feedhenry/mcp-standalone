package openshift

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
)

//ClientBuilder is a utility to help in the construction of OpenShift clients
type ClientBuilder struct {
	token, host, ns string
	incluster       bool
	insecure        bool
}

func NewClientBuilder(host, ns string, incluster, insecure bool) mobile.OSClientBuilder {
	cb := &ClientBuilder{
		host:      host,
		ns:        ns,
		incluster: incluster,
		insecure:  insecure,
	}
	return cb

}

func (cb *ClientBuilder) WithToken(token string) mobile.OSClientBuilder {
	return &ClientBuilder{
		host:      cb.host,
		token:     token,
		ns:        cb.ns,
		incluster: cb.incluster,
		insecure:  cb.insecure,
	}
}
func (cb *ClientBuilder) WithNamespace(ns string) mobile.OSClientBuilder {
	return &ClientBuilder{
		host:      cb.host,
		token:     cb.token,
		ns:        ns,
		incluster: cb.incluster,
	}
}
func (cb *ClientBuilder) WithHost(host string) mobile.OSClientBuilder {
	return &ClientBuilder{
		host:      host,
		token:     cb.token,
		ns:        cb.ns,
		incluster: cb.incluster,
		insecure:  cb.insecure,
	}
}
func (cb *ClientBuilder) WithHostAndNamespace(host, ns string) mobile.OSClientBuilder {
	return &ClientBuilder{
		host:      host,
		token:     cb.token,
		ns:        ns,
		incluster: cb.incluster,
		insecure:  cb.insecure,
	}
}
func (cb *ClientBuilder) BuildClient() (client.Interface, error) {
	//using own http client here as it will be replaced by the oc client once 3.8 arrives

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cb.insecure},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = 15 * time.Second
	return client.NewClient(cb.ns, cb.host, cb.token, httpClient), nil
}
