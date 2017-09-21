package metrics

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/clients"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/k8s"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"

	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/feedhenry/mcp-standalone/pkg/mock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func buildDefaultTestTokenClientBuilder(kclient kubernetes.Interface) mobile.TokenScopedClientBuilder {
	logger := logrus.StandardLogger()
	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	svcRepoBuilder := data.NewServiceRepoBuilder()
	svcRepoBuilder = svcRepoBuilder.WithClient(kclient.CoreV1().Secrets("test"))
	mounterBuilder := k8s.NewMounterBuilder("test")
	clientBuilder := clients.NewTokenScopedClientBuilder(cb, svcRepoBuilder, mounterBuilder, "test", logger)
	return clientBuilder
}

func TestKeycloak_Gather(t *testing.T) {

	cases := []struct {
		Name        string
		ExpectError bool
		Client      func() kubernetes.Interface
		Validate    func(t *testing.T, metrics []*metric)
		Requester   func(t *testing.T) mobile.ExternalHTTPRequester
	}{
		{
			Name: "test gather gathers at expected",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "keycloak",
								},
								Data: map[string][]byte{
									"admin_password": []byte("admin"),
									"admin_username": []byte("admin"),
									"uri":            []byte("keycloak-project2.192.168.37.1.nip.io"),
									"realm":          []byte("test"),
									"name":           []byte("keycloak"),
									"type":           []byte("keycloak"),
								},
							},
						},
					}, nil
				})
				return client
			},
			Requester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						if path == "/auth/realms/master/protocol/openid-connect/token" {
							bod := bytes.NewReader([]byte(`{"expires_in":102,"access_token":"token"}`))
							return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bod)}, nil
						} else if path == "/auth/admin/realms/test/client-session-stats" {
							bod := bytes.NewReader([]byte(`[{"clientID":"idone","active":"3"}]`))
							return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bod)}, nil
						} else if path == "/auth/admin/realms/test/events" {
							bod := bytes.NewReader([]byte(`[{"type":"LOGIN"}, {"type":"LOGIN"}, {"type":"LOGOUT"}]`))
							return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bod)}, nil
						}

						return nil, errors.New("unknown path " + path + " don't know how to respond")
					},
				}
			},
			Validate: func(t *testing.T, metrics []*metric) {
				if len(metrics) == 0 {
					t.Fatal("expected some metrics but got none")
				}
				loginFound := false
				logoutFound := false
				for _, m := range metrics {
					if m.Type == "idone" && m.YValue != 3 {
						t.Fatalf("expected the value of client logins for idone to be 3 but got %v", m.YValue)
					}
					if m.Type == "LOGOUT" {
						logoutFound = true
						if m.YValue != 1 {
							t.Fatalf("expected 1 logout but got %v ", m.YValue)
						}
					}
					if m.Type == "LOGIN" {
						loginFound = true
						if m.YValue != 2 {
							t.Fatalf("expected 2 logins but got %v ", m.YValue)
						}
					}
				}
				if !logoutFound && loginFound {
					t.Fatalf("expected to find a login metric and a log out metric but found none")
				}
			},
		},
		{
			Name:        "test error when service not available",
			ExpectError: true,
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, nil
				})
				return client
			},
			Requester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						return nil, errors.New("unknown path did not expect outbound request " + path + " don't know how to respond")
					},
				}
			},
			Validate: func(t *testing.T, metrics []*metric) {
				if metrics != nil {
					t.Fatalf("did not expect an metrics when service not available")
				}
			},
		},
		{
			Name: "test gather gathers at expected when a single call fails",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "keycloak",
								},
								Data: map[string][]byte{
									"admin_password": []byte("admin"),
									"admin_username": []byte("admin"),
									"uri":            []byte("keycloak-project2.192.168.37.1.nip.io"),
									"realm":          []byte("test"),
									"name":           []byte("keycloak"),
									"type":           []byte("keycloak"),
								},
							},
						},
					}, nil
				})
				return client
			},
			Requester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						if path == "/auth/realms/master/protocol/openid-connect/token" {
							bod := bytes.NewReader([]byte(`{"expires_in":102,"access_token":"token"}`))
							return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bod)}, nil
						} else if path == "/auth/admin/realms/test/client-session-stats" {
							bod := bytes.NewReader([]byte(`[{"clientID":"idone","active":"3"}]`))
							return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bod)}, nil
						} else if path == "/auth/admin/realms/test/events" {
							bod := bytes.NewReader([]byte(`error bad request`))
							return &http.Response{StatusCode: 400, Body: ioutil.NopCloser(bod)}, nil
						}

						return nil, errors.New("unknown path " + path + " don't know how to respond")
					},
				}
			},
			Validate: func(t *testing.T, metrics []*metric) {
				if len(metrics) != 1 {
					t.Fatal("expected some metrics but got none")
				}

				for _, m := range metrics {
					if m.Type == "idone" && m.YValue != 3 {
						t.Fatalf("expected the value of client logins for idone to be 3 but got %v", m.YValue)
					}
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			httpClientBuilder := &mock.HttpClientBuilder{Requester: tc.Requester(t)}
			kc := NewKeycloak(httpClientBuilder, buildDefaultTestTokenClientBuilder(tc.Client()), logrus.StandardLogger())
			metrics, err := kc.Gather()
			if err == nil && tc.ExpectError {
				t.Fatal("expected an error but got none")
			}
			if err != nil && !tc.ExpectError {
				t.Fatalf("did not expect an error but got one %s", err)
			}
			if tc.Validate != nil {
				tc.Validate(t, metrics)
			}
		})
	}
}
