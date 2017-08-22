package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/data"
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/feedhenry/mobile-server/pkg/mobile/client"
	"github.com/feedhenry/mobile-server/pkg/mock"
	"github.com/feedhenry/mobile-server/pkg/web"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func setupMobileAppHandler(kclient kubernetes.Interface, cruder corev1.ConfigMapInterface, svcCruder corev1.SecretInterface) http.Handler {
	r := web.NewRouter()
	logger := logrus.New()

	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	appRepoBuilder := data.NewMobileAppRepoBuilder()
	appRepoBuilder = appRepoBuilder.WithClient(cruder)
	svcRepoBuilder := data.NewServiceRepoBuilder()
	svcRepoBuilder = svcRepoBuilder.WithClient(svcCruder)
	clientBuilder := client.NewTokenScopedClientBuilder(cb, appRepoBuilder, svcRepoBuilder, "test", logger)
	handler := web.NewMobileAppHandler(logger, clientBuilder)
	web.MobileAppRoute(r, handler)
	return web.BuildHTTPHandler(r, nil)
}

func TestReadMobileApp(t *testing.T) {
	cases := []struct {
		Name       string
		Clients    func() (kubernetes.Interface, corev1.ConfigMapInterface)
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
			Clients: func() (kubernetes.Interface, corev1.ConfigMapInterface) {
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
				return client, client.CoreV1().ConfigMaps("test")
			},
			StatusCode: 200,
		},
		{
			Name: "test read 404 via api when no app",
			App: &mobile.App{
				Name: "app",
			},
			Clients: func() (kubernetes.Interface, corev1.ConfigMapInterface) {
				client := &fake.Clientset{}
				client.AddReactor("get", "configmaps", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, nil, &kerror.StatusError{ErrStatus: metav1.Status{Code: 404}}
				})
				return client, client.CoreV1().ConfigMaps("test")
			},
			StatusCode: 404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			kclient, configmapClient := tc.Clients()
			secretClient := kclient.CoreV1().Secrets("test")
			handler := setupMobileAppHandler(kclient, configmapClient, secretClient)
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
