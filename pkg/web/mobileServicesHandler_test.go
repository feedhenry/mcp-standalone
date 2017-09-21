package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/k8s"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/clients"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/metrics"
	"github.com/feedhenry/mcp-standalone/pkg/mock"
	"github.com/feedhenry/mcp-standalone/pkg/web"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	ktesting "k8s.io/client-go/testing"
)

func setupMobileServiceHandler(kclient kubernetes.Interface) http.Handler {
	r := web.NewRouter()
	logger := logrus.StandardLogger()
	if nil == kclient {
		kclient = &fake.Clientset{}
	}
	clientBuilder := buildDefaultTestTokenClientBuilder(kclient)
	metricGetter := &metrics.MetricsService{}
	ms := &integration.MobileService{}
	handler := web.NewMobileServiceHandler(logger, ms, clientBuilder, metricGetter)
	web.MobileServiceRoute(r, handler)
	return web.BuildHTTPHandler(r, nil)
}

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

func TestListMobileServices(t *testing.T) {
	cases := []struct {
		Name       string
		Client     func() kubernetes.Interface
		StatusCode int
		Validate   func(svs []*mobile.Service, t *testing.T)
	}{
		{
			Name: "test list mobile services ok",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"name": []byte("fh-sync-server"),
									"type": []byte("fh-sync-server"),
								},
							},
						}}, nil
				})
				return client
			},
			Validate: func(svs []*mobile.Service, t *testing.T) {
				if nil == svs {
					t.Fatal("expected some mobile services but got none")
				}
				if len(svs) != 1 {
					t.Fatalf("expected 1 mobile service but got %d ", len(svs))
				}
				s := svs[0]
				if s.Name != "fh-sync-server" {
					t.Fatal("expected the mobile service name to be fh-sync-server")
				}
			},
			StatusCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileServiceHandler(tc.Client())
			server := httptest.NewServer(handler)
			defer server.Close()
			res, err := http.Get(server.URL + "/mobileservice")
			if err != nil {
				t.Fatal("did not expect an error requesting mobile services ", err)
			}
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("expected a response code of %d but got %d ", tc.StatusCode, res.StatusCode)
			}
			if nil != tc.Validate {
				decoder := json.NewDecoder(res.Body)
				svs := []*mobile.Service{}
				if err := decoder.Decode(&svs); err != nil {
					t.Fatal("unexpected error decoding the services response ", err)
				}
				tc.Validate(svs, t)
			}
		})
	}
}

func TestConfigure(t *testing.T) {
	cases := []struct {
		Name       string
		Client     func() kubernetes.Interface
		StatusCode int
		Validate   func(r *http.Response, t *testing.T)
	}{
		{
			Name: "test configuring fh-sync-server for keycloak",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "fh-sync-server-secret",
								},
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"type": []byte("fh-sync-server"),
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "keycloak-public-client",
								},
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"type": []byte("keycloak"),
								},
							},
						}}, nil
				})
				client.AddReactor("get", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name: "keycloak-public-client",
						},
						Data: map[string][]byte{},
					}, nil
				})
				client.AddReactor("get", "deployments", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1beta1.Deployment{
						Spec: v1beta1.DeploymentSpec{
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									Volumes: []v1.Volume{},
									Containers: []v1.Container{
										{
											Name:         "fh-sync-server",
											VolumeMounts: []v1.VolumeMount{},
										},
									},
								},
							},
						},
					}, nil
				})
				client.AddReactor("Update", "deployments", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1beta1.Deployment{
						Spec: v1beta1.DeploymentSpec{
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									Volumes: []v1.Volume{},
									Containers: []v1.Container{
										{
											Name:         "fh-sync-server",
											VolumeMounts: []v1.VolumeMount{},
										},
									},
								},
							},
						},
					}, nil
				})
				return client
			},
			StatusCode: http.StatusOK,
			Validate: func(r *http.Response, t *testing.T) {
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				if len(bodyBytes) != 0 {
					t.Fatalf("expected zero bytes in response body, got %d", len(bodyBytes))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileServiceHandler(tc.Client())
			server := httptest.NewServer(handler)
			defer server.Close()
			res, err := http.Post(server.URL+"/mobileservice/configure/fh-sync-server/keycloak-public-client", "text/plain", strings.NewReader(""))
			if err != nil {
				t.Fatal("did not expect an error requesting mobile services ", err)
			}
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("expected a response code of %d but got %d ", tc.StatusCode, res.StatusCode)
			}
			if nil != tc.Validate {
				tc.Validate(res, t)
			}
		})
	}
}

func TestDeconfigure(t *testing.T) {
	cases := []struct {
		Name       string
		Client     func() kubernetes.Interface
		StatusCode int
		Validate   func(r *http.Response, t *testing.T)
	}{
		{
			Name: "test deconfiguring fh-sync-server for keycloak",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "fh-sync-server",
								},
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"type": []byte("fh-sync-server"),
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "keycloak-secret",
								},
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"name": []byte("keycloak"),
									"type": []byte("keycloak"),
								},
							},
						}}, nil
				})
				client.AddReactor("get", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name: "keycloak-public-client",
						},
						Data: map[string][]byte{},
					}, nil
				})
				client.AddReactor("get", "deployments", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1beta1.Deployment{
						Spec: v1beta1.DeploymentSpec{
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									Volumes: []v1.Volume{
										{
											Name: "keycloak-public-client",
										},
									},
									Containers: []v1.Container{
										{
											Name: "fh-sync-server",
											VolumeMounts: []v1.VolumeMount{
												{
													Name: "keycloak-public-client",
												},
											},
										},
									},
								},
							},
						},
					}, nil
				})
				client.AddReactor("Update", "deployments", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1beta1.Deployment{
						Spec: v1beta1.DeploymentSpec{
							Template: v1.PodTemplateSpec{
								Spec: v1.PodSpec{
									Volumes: []v1.Volume{},
									Containers: []v1.Container{
										{
											Name:         "fh-sync-server",
											VolumeMounts: []v1.VolumeMount{},
										},
									},
								},
							},
						},
					}, nil
				})
				return client
			},
			StatusCode: http.StatusOK,
			Validate: func(r *http.Response, t *testing.T) {
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				if len(bodyBytes) != 0 {
					t.Fatalf("expected zero bytes in response body, got %d", len(bodyBytes))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileServiceHandler(tc.Client())
			server := httptest.NewServer(handler)
			defer server.Close()
			req, err := http.NewRequest("DELETE", server.URL+"/mobileservice/configure/fh-sync-server/keycloak-public-client", strings.NewReader(""))
			if err != nil {
				t.Fatal("did not expect an error creating a http requets", err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("did not expect an error requesting mobile services ", err)
			}
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("expected a response code of %d but got %d ", tc.StatusCode, res.StatusCode)
			}
			if nil != tc.Validate {
				tc.Validate(res, t)
			}
		})
	}
}

func TestCreateMobileService(t *testing.T) {
	cases := []struct {
		Name          string
		ExpectError   bool
		StatusCode    int
		MobileService *mobile.Service
	}{
		{
			Name:       "test create mobile service is ok",
			StatusCode: 201,
			MobileService: &mobile.Service{
				Name: "mykeycloak",
				Host: "https://sdasdd.com",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileServiceHandler(nil)
			server := httptest.NewServer(handler)
			defer server.Close()
			bod, err := json.Marshal(tc.MobileService)
			if err != nil {
				t.Fatalf("failed to marshal body for service create request %v", err)
			}
			res, err := http.Post(server.URL+"/mobileservice", "application/json", bytes.NewReader(bod))
			if err != nil {
				t.Fatal("did not expect an error creating a mobile services ", err)
			}
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("expected a response code of %d but got %d ", tc.StatusCode, res.StatusCode)
			}
		})
	}
}

func TestMobileServiceHandler_Delete(t *testing.T) {
	cases := []struct {
		Name        string
		StatusCode  int
		ServiceName string
		ExpectError bool
	}{
		{
			Name:        "test delete mobile service ok",
			StatusCode:  200,
			ServiceName: "test-service",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileServiceHandler(nil)
			server := httptest.NewServer(handler)
			defer server.Close()

			req, err := http.NewRequest("DELETE", server.URL+"/mobileservice/"+tc.ServiceName, nil)
			if err != nil {
				t.Fatal("did not expect an error creating delete request ", err)
			}
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("did not expect an error doing the delete request ", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.StatusCode {
				t.Fatalf("expected response code %v but got %v ", http.StatusOK, res.StatusCode)
			}
		})
	}

}
