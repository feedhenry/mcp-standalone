package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
)

type HttpClientBuilder struct {
	insecure bool
	timeout  int
}

func NewHttpClientBuilder() *HttpClientBuilder {
	return &HttpClientBuilder{
		insecure: false,
		timeout:  30,
	}
}

func (hcb *HttpClientBuilder) Insecure(i bool) *HttpClientBuilder {
	hcb.insecure = i
	return hcb
}

func (hcb *HttpClientBuilder) Timeout(t int) *HttpClientBuilder {
	hcb.timeout = t
	return hcb
}

func (hcb *HttpClientBuilder) Build() mobile.ExternalHTTPRequester {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: hcb.insecure},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = time.Second * time.Duration(hcb.timeout)
	return client
}
