package metrics

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/clients"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/k8s"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"

	"github.com/feedhenry/mcp-standalone/pkg/mock"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func buildDefaultTestTokenClientBuilder(kclient kubernetes.Interface) mobile.TokenScopedClientBuilder {
	logger := logrus.StandardLogger()
	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	appRepoBuilder := data.NewMobileAppRepoBuilder()
	appRepoBuilder = appRepoBuilder.WithClient(kclient.CoreV1().ConfigMaps("test"))
	svcRepoBuilder := data.NewServiceRepoBuilder()
	svcRepoBuilder = svcRepoBuilder.WithClient(kclient.CoreV1().Secrets("test"))
	mounterBuilder := k8s.NewMounterBuilder("test")
	clientBuilder := clients.NewTokenScopedClientBuilder(cb, appRepoBuilder, svcRepoBuilder, mounterBuilder, "test", logger)
	return clientBuilder
}

func TestKeycloak_Gather(t *testing.T) {
	cases := []struct {
		Name        string
		ExpectError bool
		Client      func() kubernetes.Interface
		Validate    func(metrics []*metric)
	}{
		{
			Name: "test gather gathers at expected",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				return client
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mr := mock.Requester{}

			kc := NewKeycloak(mr, buildDefaultTestTokenClientBuilder(tc.Client()), "keycloak", logrus.StandardLogger())
			metrics, err := kc.Gather()

		})
	}
}
