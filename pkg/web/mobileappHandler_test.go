package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/mock"
	"github.com/feedhenry/mcp-standalone/pkg/web"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func setupMobileAppHandler(kclient kubernetes.Interface) http.Handler {
	r := web.NewRouter()
	logger := logrus.StandardLogger()
	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	appRepoBuilder := data.NewMobileAppRepoBuilder(cb, "test", "test")
	appService := &app.Service{}
	handler := web.NewMobileAppHandler(logger, appRepoBuilder, appService)
	web.MobileAppRoute(r, handler)
	return web.BuildHTTPHandler(r, nil)
}

func TestCreateMobileApp(t *testing.T) {
	cases := []struct {
		Name       string
		Client     func() kubernetes.Interface
		StatusCode int
		Validate   func(app *mobile.App, t *testing.T)
		MobileApp  *mobile.App
	}{
		{
			Name: "test create via api is ok",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("get", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.ConfigMap{
						Data: map[string]string{},
					}, nil
				})
				return client
			},
			StatusCode: 201,
			MobileApp: &mobile.App{
				Name:       "myApp",
				ClientType: "cordova",
			},
		},
		{
			Name: "test create via api errors",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("get", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.ConfigMap{
						Data: map[string]string{},
					}, nil
				})
				return client
			},
			StatusCode: 400,
			MobileApp: &mobile.App{
				Name: "myApp",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			kclient := tc.Client()
			handler := setupMobileAppHandler(kclient)
			server := httptest.NewServer(handler)
			defer server.Close()
			body, err := json.Marshal(tc.MobileApp)
			if err != nil {
				t.Fatal("failed to marshal body for app, create request", err)
			}
			res, err := http.Post(server.URL+"/mobileapp", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Fatal("encountered error performed create request", err)
			}
			defer res.Body.Close()
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("incorrect status code %d, expected %d", res.StatusCode, tc.StatusCode)
			}
		})
	}
}

func TestReadMobileApp(t *testing.T) {
	cases := []struct {
		Name       string
		Clients    func() kubernetes.Interface
		App        *mobile.App
		StatusCode int
		Validate   func(app *mobile.App, t *testing.T)
	}{
		{
			Name: "test read happpy via api",
			Validate: func(app *mobile.App, t *testing.T) {
				if app == nil {
					t.Fatal("expectd an app but it was nil")
				}
				if app.Name != "app" {
					t.Fatalf("expected the app name to be app but got %s ", app.Name)
				}
				if app.ClientType != "iOS" {
					t.Fatalf("expected the app clientType to be android but got %s ", app.ClientType)
				}
			},
			App: &mobile.App{
				Name: "app",
			},
			Clients: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("get", "configmaps", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string]string{
							"name":       "app",
							"clientType": "iOS",
						},
					}, nil
				})
				return client
			},
			StatusCode: 200,
		},
		{
			Name: "test read 404 via api when no app",
			App: &mobile.App{
				Name: "app",
			},
			Clients: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("get", "configmaps", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, nil, &kerror.StatusError{ErrStatus: metav1.Status{Code: 404}}
				})
				return client
			},
			StatusCode: 404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			kclient := tc.Clients()
			handler := setupMobileAppHandler(kclient)
			server := httptest.NewServer(handler)
			defer server.Close()
			req, err := http.NewRequest("GET", server.URL+"/mobileapp/"+tc.App.Name, nil)
			if err != nil {
				t.Fatalf("did not expect an error creating request %v", err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("did not expect an error doing request %v", err)
			}
			if tc.StatusCode != res.StatusCode {
				t.Fatalf("expected status code %v but got %v ", tc.StatusCode, res.StatusCode)
			}
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				decoder := json.NewDecoder(res.Body)
				app := &mobile.App{}
				if err := decoder.Decode(app); err != nil {
					t.Fatalf("failed to decode reponse body into mobile app %v", err)
				}
				tc.Validate(app, t)
			}
		})
	}
}
