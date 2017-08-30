package web_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/client"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
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
	clientBuilder := buildDefaultTestTokenClientBuilder(kclient)
	ms := &integration.MobileService{}
	handler := web.NewMobileServiceHandler(logger, ms, clientBuilder)
	web.MobileServiceRoute(r, handler)
	return web.BuildHTTPHandler(r, nil, nil)
}

func buildDefaultTestTokenClientBuilder(kclient kubernetes.Interface) mobile.TokenScopedClientBuilder {
	logger := logrus.StandardLogger()
	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	appRepoBuilder := data.NewMobileAppRepoBuilder()
	appRepoBuilder = appRepoBuilder.WithClient(kclient.CoreV1().ConfigMaps("test"))
	svcRepoBuilder := data.NewServiceRepoBuilder()
	svcRepoBuilder = svcRepoBuilder.WithClient(kclient.CoreV1().Secrets("test"))
	clientBuilder := client.NewTokenScopedClientBuilder(cb, appRepoBuilder, svcRepoBuilder, "test", logger)
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
									"name": []byte("fh-sync-server"),
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "keycloak-secret",
								},
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"name": []byte("keycloak"),
								},
							},
						}}, nil
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
				return client
			},
			StatusCode: http.StatusOK,
			Validate: func(r *http.Response, t *testing.T) {
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				res := &v1beta1.Deployment{}
				err := json.Unmarshal(bodyBytes, res)
				if err != nil {
					t.Fatal(err)
				}

				if len(res.Spec.Template.Spec.Volumes) == 0 {
					t.Fatal("Expected a volume to be added to the deployment")
				}
				if len(res.Spec.Template.Spec.Containers[0].VolumeMounts) == 0 {
					t.Fatal("Expected a volumemount to be added to the container in the deployment")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileServiceHandler(tc.Client())
			server := httptest.NewServer(handler)
			defer server.Close()
			res, err := http.Post(server.URL+"/mobileservice/configure", "text/plain", strings.NewReader("{\"component\": \"fh-sync-server\", \"service\": \"keycloak\", \"namespace\": \"myproject\"}"))
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
