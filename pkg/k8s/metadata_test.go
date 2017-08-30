package k8s

import (
	"errors"
	"net/http"
	"testing"

	"bytes"
	"io/ioutil"

	"fmt"

	"github.com/feedhenry/mcp-standalone/pkg/mock"
)

func TestGetMetadata(t *testing.T) {
	cases := []struct {
		Name            string
		K8shost         string
		RequestResponse func(host string, path string, method string, t *testing.T) (*http.Response, error)
		ExpectError     bool
		Validate        func(metadata *Metadata, err error, t *testing.T)
	}{
		{
			Name:    "test get metadata ok",
			K8shost: "https://127.0.0.1:443",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				buf := bytes.NewBufferString(`{
				  "issuer": "https://mockk8shost.example.com",
				  "authorization_endpoint": "https://mockk8shost.example.com/oauth/authorize",
				  "token_endpoint": "https://mockk8shost.example.com/oauth/token",
				  "scopes_supported": [
					"user:info",
					"user:check-access"
				  ],
				  "response_types_supported": [
					"code",
					"token"
				  ],
				  "grant_types_supported": [
					"authorization_code",
					"implicit"
				  ],
				  "code_challenge_methods_supported": [
					"plain",
					"S256"
				  ]
				}`)
				if host == "127.0.0.1:443" && path == "/.well-known/oauth-authorization-server" {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(buf),
					}, nil
				}
				return nil, errors.New(fmt.Sprintf("unexpected host or path %s %s ", host, path))
			},
			Validate: func(metadata *Metadata, err error, t *testing.T) {
				if metadata == nil {
					t.Fatal("expected metadata but got none")
				}
				expected := "https://mockk8shost.example.com"
				if metadata.Issuer != expected {
					t.Fatalf("expected the metadata host to be '%s' but got '%s'", expected, metadata.Issuer)
				}
			},
		},
		{
			Name:    "test get metadata scope missing",
			K8shost: "https://127.0.0.1:443",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				buf := bytes.NewBufferString(`{
				  "issuer": "https://mockk8shost.example.com",
				  "authorization_endpoint": "https://mockk8shost.example.com/oauth/authorize",
				  "token_endpoint": "https://mockk8shost.example.com/oauth/token",
				  "scopes_supported": [
					"user:info"
				  ],
				  "response_types_supported": [
					"code",
					"token"
				  ],
				  "grant_types_supported": [
					"authorization_code",
					"implicit"
				  ],
				  "code_challenge_methods_supported": [
					"plain",
					"S256"
				  ]
				}`)
				if host == "127.0.0.1:443" && path == "/.well-known/oauth-authorization-server" {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(buf),
					}, nil
				}
				return nil, errors.New(fmt.Sprintf("unexpected host or path %s %s ", host, path))
			},
			ExpectError: true,
			Validate: func(metadata *Metadata, err error, t *testing.T) {
				expected := "OpenShift server metadata is missing required scopes. Got [user:info]"
				if err.Error() != expected {
					t.Fatal(fmt.Sprintf("expected error %s but got %s", expected, err))
				}
			},
		},
		{
			Name:    "test get metadata response type missing",
			K8shost: "https://127.0.0.1:443",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				buf := bytes.NewBufferString(`{
				  "issuer": "https://mockk8shost.example.com",
				  "authorization_endpoint": "https://mockk8shost.example.com/oauth/authorize",
				  "token_endpoint": "https://mockk8shost.example.com/oauth/token",
				  "scopes_supported": [
					"user:info",
					"user:check-access"
				  ],
				  "response_types_supported": [
					"code"
				  ],
				  "grant_types_supported": [
					"authorization_code",
					"implicit"
				  ],
				  "code_challenge_methods_supported": [
					"plain",
					"S256"
				  ]
				}`)
				if host == "127.0.0.1:443" && path == "/.well-known/oauth-authorization-server" {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(buf),
					}, nil
				}
				return nil, errors.New(fmt.Sprintf("unexpected host or path %s %s ", host, path))
			},
			ExpectError: true,
			Validate: func(metadata *Metadata, err error, t *testing.T) {
				expected := "OpenShift server metadata is missing required response types. Got [code]"
				if err.Error() != expected {
					t.Fatal(fmt.Sprintf("expected error %s but got %s", expected, err))
				}
			},
		},
		{
			Name:    "test get metadata grant type missing",
			K8shost: "https://127.0.0.1:443",
			RequestResponse: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
				buf := bytes.NewBufferString(`{
				  "issuer": "https://mockk8shost.example.com",
				  "authorization_endpoint": "https://mockk8shost.example.com/oauth/authorize",
				  "token_endpoint": "https://mockk8shost.example.com/oauth/token",
				  "scopes_supported": [
					"user:info",
					"user:check-access"
				  ],
				  "response_types_supported": [
					"code",
					"token"
				  ],
				  "grant_types_supported": [
					"implicit"
				  ],
				  "code_challenge_methods_supported": [
					"plain",
					"S256"
				  ]
				}`)
				if host == "127.0.0.1:443" && path == "/.well-known/oauth-authorization-server" {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(buf),
					}, nil
				}
				return nil, errors.New(fmt.Sprintf("unexpected host or path %s %s ", host, path))
			},
			ExpectError: true,
			Validate: func(metadata *Metadata, err error, t *testing.T) {
				expected := "OpenShift server metadata is missing required grant types. Got [implicit]"
				if err.Error() != expected {
					t.Fatal(fmt.Sprintf("expected error %s but got %s", expected, err))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			metadata, err := GetMetadata(tc.K8shost, &mock.Requester{Responder: tc.RequestResponse})
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %v ", err)
			}
			tc.Validate(metadata, err, t)
		})
	}
}
