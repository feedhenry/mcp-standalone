package mock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
)

type Requester struct {
	Test      *testing.T
	Responder func(host string, path string, method string, t *testing.T) (*http.Response, error)
}

func (mr *Requester) Do(req *http.Request) (*http.Response, error) {
	return mr.Responder(req.Host, req.URL.Path, req.Method, mr.Test)
}

func (mr *Requester) Get(fullUrl string) (*http.Response, error) {
	parsedUrl, _ := url.Parse(fullUrl)

	return mr.Responder(parsedUrl.Host, parsedUrl.Path, "GET", mr.Test)
}

// HttpClientBuilder mocks out the clients.ExternalRequester that is a mobile.ExternalHTTPRequester
type HttpClientBuilder struct {
	Requester mobile.ExternalHTTPRequester
}

func (hcb *HttpClientBuilder) Insecure(i bool) mobile.HTTPRequesterBuilder {
	return hcb
}
func (hcb *HttpClientBuilder) Timeout(t int) mobile.HTTPRequesterBuilder {
	return hcb
}
func (hcb *HttpClientBuilder) Build() mobile.ExternalHTTPRequester {
	return hcb.Requester
}
