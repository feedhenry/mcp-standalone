package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	v1 "k8s.io/client-go/pkg/api/v1"
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
									"name": []byte("fh-sync"),
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
				if s.Name != "fh-sync" {
					t.Fatal("expected the mobile service name to be fh-sync")
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
