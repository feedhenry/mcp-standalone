package integration_test

import (
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestSDKService_GenerateMobileServiceConfigs(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() mobile.ServiceCruder
		Validate    func(sdkConfigs map[string]*mobile.ServiceConfig, t *testing.T)
		ExpectError bool
	}{
		{
			Name: "test generate service configs ok",
			Client: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(a ktesting.Action) (bool, runtime.Object, error) {
					list := v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Labels: map[string]string{"group": "mobileapp"},
								},
								Data: map[string][]byte{
									"uri":  []byte("http://fh-sync.com"),
									"name": []byte("fh-sync-server"),
									"type": []byte("fh-sync-server"),
								},
							},
							{
								Data: map[string][]byte{
									"uri": []byte("http://somehost.com"),
								},
							},
						},
					}
					return true, &list, nil
				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
			Validate: func(sdkConfigs map[string]*mobile.ServiceConfig, t *testing.T) {
				if nil == sdkConfigs {
					t.Fatal("expected sdk configs but got none")
				}
				if v, ok := sdkConfigs["fh-sync-server"]; ok {
					configValues := v.Config.(map[string]string)
					if _, ok := configValues["uri"]; !ok {
						t.Fatalf("expected a uri in the service config")
					}
				} else {
					t.Fatal("expected fh-sync-server to be in the config response")
				}
			},
		},
		{
			Name: "test generate service config fails on error",
			Client: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(a ktesting.Action) (bool, runtime.Object, error) {
					return true, nil, errors.New("unexpected error")
				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
			ExpectError: true,
			Validate: func(sdkConfigs map[string]*mobile.ServiceConfig, t *testing.T) {
				if nil != sdkConfigs {
					t.Fatalf("did not expect sdkConfigs but got %v", sdkConfigs)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			sdkService := &integration.SDKService{}
			sdkConfig, err := sdkService.GenerateMobileServiceConfigs(tc.Client())
			if err == nil && tc.ExpectError {
				t.Fatal("expected an error but got none")
			}
			if err != nil && !tc.ExpectError {
				t.Fatalf("did not expect an error but got one %s", err)
			}
			tc.Validate(sdkConfig, t)
		})
	}

}
