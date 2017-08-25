package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/web"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func setupSDKHandler(kclient kubernetes.Interface) http.Handler {
	r := web.NewRouter()
	logger := logrus.StandardLogger()
	clientBuilder := buildDefaultTestTokenClientBuilder(kclient)
	ms := &integration.MobileService{}
	sdkConfigHandler := web.NewSDKConfigHandler(logger, ms, clientBuilder)
	web.SDKConfigRoute(r, sdkConfigHandler)
	return web.BuildHTTPHandler(r, nil)
}

func TestSDKConifg(t *testing.T) {
	var apiKey = "supersecure"
	cases := []struct {
		Name       string
		Client     func() kubernetes.Interface
		StatusCode int
	}{
		{
			StatusCode: http.StatusOK,
			Name:       "test get sdk config ok",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("get", "configmaps", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string]string{
							"name":       "app",
							"clientType": "iOS",
							"apiKey":     apiKey,
						},
					}, nil
				})
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
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client()
			server := httptest.NewServer(setupSDKHandler(client))
			defer server.Close()
			req, err := http.NewRequest("GET", server.URL+"/sdk/mobileapp/app/config", nil)
			if err != nil {
				t.Fatal("did not expect error setting up the request ", err)
			}
			req.Header.Add(mobile.AppAPIKeyHeader, apiKey)
			httpclient := &http.Client{}
			httpclient.Timeout = 2 * time.Second
			res, err := httpclient.Do(req)
			if err != nil {
				t.Fatal("did not expect an error making config request ", err)
			}
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("expected status code %v but got %v", tc.StatusCode, res.StatusCode)
			}
			if tc.StatusCode == http.StatusOK {
				decoder := json.NewDecoder(res.Body)
				data := map[string]*mobile.Service{}
				if err := decoder.Decode(&data); err != nil {
					t.Fatal("did not expect an error decoding the config response ", err)
				}
				if _, ok := data["fh-sync"]; !ok {
					t.Fatal("expected fh-sync to be in the config response but it was missing ")
				}
			}
		})
	}
}
