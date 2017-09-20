package integration_test

import (
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

type mockAuthChecker struct {
	checkBoolRes bool
	checkErrRes  error
}

func (mac *mockAuthChecker) Check(resource, namespace string) (bool, error) {
	return mac.checkBoolRes, mac.checkErrRes
}

func TestMobileServiceDiscovery(t *testing.T) {
	cases := []struct {
		Name          string
		ExpectError   bool
		ServiceCruder func() mobile.ServiceCruder
		authChecker   func() mobile.AuthChecker
		Validate      func(svs []*mobile.Service, t *testing.T)
	}{
		{
			Name: "test discover mobile services ok",
			authChecker: func() mobile.AuthChecker {
				return &mockAuthChecker{}
			},
			ServiceCruder: func() mobile.ServiceCruder {
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
			Validate: func(svs []*mobile.Service, t *testing.T) {
				if svs == nil {
					t.Fatalf("expected a list of services but got none")
				}
				if len(svs) != 1 {
					t.Fatalf("expected 1 service to be returned but got %v", len(svs))
				}
				svc := svs[0]
				if svc.Host != "http://fh-sync.com" {
					t.Fatalf("unexpected host %s", svc.Host)
				}
				if svc.Name != "fh-sync-server" {
					t.Fatalf("expected svc name to be fh-sync-server but got %s", svc.Name)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			is := integration.MobileService{}
			svcs, err := is.DiscoverMobileServices(tc.ServiceCruder(), tc.authChecker())
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none!")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect and err but got one %v", err)
			}
			tc.Validate(svcs, t)
		})
	}
}
