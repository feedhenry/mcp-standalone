package jenkins_test

import (
	"testing"

	"net/http"

	"bytes"

	"io/ioutil"

	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/jenkins"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mock"
)

func TestClientRetrieve(t *testing.T) {
	cases := []struct {
		Name              string
		ExternalRequester func(t *testing.T) mobile.ExternalHTTPRequester
		Location          *url.URL
		ExpectError       bool
	}{
		{
			Name: "test retrieve artifact ok",
			ExternalRequester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						var res = &http.Response{
							StatusCode: 200,
						}
						if path == "/job/localmcp-androidapp1/1/artifact/platforms/android/build/outputs/apk/android-debug.apk" {
							res.Body = ioutil.NopCloser(bytes.NewBufferString("testdata"))
						} else {
							res.Body = ioutil.NopCloser(bytes.NewBufferString(`[{"id":"n1","name":"android-debug.apk","path":"platforms/android/build/outputs/apk/android-debug.apk","url":"/job/localmcp-androidapp1/1/artifact/platforms/android/build/outputs/apk/android-debug.apk","size":4361695}]`))
						}
						return res, nil
					},
				}
			},
			Location: &url.URL{
				Scheme: "http",
				Host:   "jenkins.com",
				Path:   "/job/localmcp-androidapp1/1/wfapi/artifacts",
			},
		},
		{
			Name: "test fails when no artifact",
			ExternalRequester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						var res = &http.Response{
							StatusCode: 200,
						}
						res.Body = ioutil.NopCloser(bytes.NewBufferString(`[]`))

						return res, nil
					},
				}
			},
			ExpectError: true,
			Location: &url.URL{
				Scheme: "http",
				Host:   "jenkins.com",
				Path:   "/job/localmcp-androidapp1/1/wfapi/artifacts",
			},
		},
		{
			Name: "test fails when non 200",
			ExternalRequester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						var res = &http.Response{
							StatusCode: 401,
						}

						res.Body = ioutil.NopCloser(bytes.NewBufferString(`unauthorised`))
						return res, nil
					},
				}
			},
			ExpectError: true,
			Location: &url.URL{
				Scheme: "http",
				Host:   "jenkins.com",
				Path:   "/job/localmcp-androidapp1/1/wfapi/artifacts",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			jc := jenkins.NewClient(tc.ExternalRequester(t), logrus.StandardLogger())
			rc, err := jc.Retrieve(tc.Location, "token")
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %v ", err)
			}
			if !tc.ExpectError {
				data, _ := ioutil.ReadAll(rc)
				if string(data) != "testdata" {
					t.Fatal("expected to get test data but got ", string(data))
				}
			}
		})
	}
}
