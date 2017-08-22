package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/data"
	"github.com/feedhenry/mobile-server/pkg/mobile/client"
	"github.com/feedhenry/mobile-server/pkg/mobile/integration"
	"github.com/feedhenry/mobile-server/pkg/mock"
	"github.com/feedhenry/mobile-server/pkg/web"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func setupSDKHandler(kclient kubernetes.Interface) http.Handler {
	r := web.NewRouter()
	logger := logrus.New()

	// TODO it is likely most handlers will need the client builder so look at potential to abstract
	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	appRepoBuilder := data.NewMobileAppRepoBuilder()
	appRepoBuilder = appRepoBuilder.WithClient(kclient.CoreV1().ConfigMaps("test"))
	svcRepoBuilder := data.NewServiceRepoBuilder()
	svcRepoBuilder = svcRepoBuilder.WithClient(kclient.CoreV1().Secrets("test"))
	clientBuilder := client.NewTokenScopedClientBuilder(cb, appRepoBuilder, svcRepoBuilder, "test", logger)
	ms := &integration.MobileService{}
	sdkConfigHandler := web.NewSDKConfigHandler(logger, ms, clientBuilder)
	web.SDKConfigRoute(r, sdkConfigHandler)
	return web.BuildHTTPHandler(r, nil)
}

func TestSDKConifg(t *testing.T) {
	cases := []struct {
		Name       string
		Client     func() kubernetes.Interface
		StatusCode int
	}{
		{
			Name: "test get sdk config ok",
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
						},
					}, nil
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
		})
	}
}
