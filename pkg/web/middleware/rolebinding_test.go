package middleware

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mock"
)

type mockHTTPClientBuilder struct {
	requester *mock.Requester
}

func (mcb *mockHTTPClientBuilder) Insecure(i bool) mobile.HTTPRequesterBuilder {
	return mcb
}

func (mcb *mockHTTPClientBuilder) Timeout(t int) mobile.HTTPRequesterBuilder {
	return mcb
}

func (mcb *mockHTTPClientBuilder) Build() mobile.ExternalHTTPRequester {
	return mcb.requester
}

func TestRolbindingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	cases := []struct {
		Name            string
		ExpectError     bool
		RequestResponse func(host string, path string, method string, t *testing.T) (*http.Response, error)
	}{
		{
			Name: "test rolebinding ok when doesn't exist",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				if path == "/oapi/v1/namespaces/test/rolebindings/edit" {
					return &http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(&buf),
					}, nil
				}
				if path == "/oapi/v1/namespaces/test/rolebindings" {
					return &http.Response{
						StatusCode: 201,
						Body:       ioutil.NopCloser(&buf),
					}, nil
				}
				return nil, errors.New("unexpected path " + path)
			},
		},
		{
			Name: "test rolebinding ok when rolebinding already exists",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				if path == "/oapi/v1/namespaces/test/rolebindings/edit" {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(&buf),
					}, nil
				}
				if path == "/oapi/v1/namespaces/test/rolebindings" {
					t.Fatal("did not expect rolbinding create to be called")
				}
				return nil, errors.New("unexpected path " + path)
			},
		},
		{
			Name: "test rolebinding fails when unauthorised",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				if path == "/oapi/v1/namespaces/test/rolebindings/edit" {
					return &http.Response{
						StatusCode: 401,
						Body:       ioutil.NopCloser(&buf),
					}, nil
				}
				if path == "/oapi/v1/namespaces/test/rolebindings" {
					t.Fatal("did not expect rolbinding create to be called")
				}
				return nil, errors.New("unexpected path " + path)
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			rb := &RoleBinding{
				clientBuilder: &mockHTTPClientBuilder{requester: &mock.Requester{Responder: tc.RequestResponse}},
				namespace:     "test",
				khost:         "http://k8s.io",
				logger:        logrus.StandardLogger(),
				Mutex:         &sync.Mutex{},
			}

			err := rb.createRoleBindingIfNotPresent("sometoken")
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one %v", err)
			}
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
		})
	}
}
