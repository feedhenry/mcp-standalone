package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
)

type ClientBuilder struct {
	insecure bool
	timeout  int
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		insecure: false,
		timeout:  30,
	}
}

func (hcb *ClientBuilder) Insecure(i bool) mobile.HTTPRequesterBuilder {
	hcb.insecure = i
	return hcb
}

func (hcb *ClientBuilder) Timeout(t int) mobile.HTTPRequesterBuilder {
	hcb.timeout = t
	return hcb
}

func (hcb *ClientBuilder) Build() mobile.ExternalHTTPRequester {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: hcb.insecure},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = time.Second * time.Duration(hcb.timeout)
	return client
}
