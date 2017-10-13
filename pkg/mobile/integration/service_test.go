package integration_test

import (
	"testing"

	"io/ioutil"
	"net/http"
	"strings"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/mock"
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

func (mac *mockAuthChecker) Check(resource, namespace string, client mobile.ExternalHTTPRequester) (bool, error) {
	return mac.checkBoolRes, mac.checkErrRes
}

func TestMobileServiceDiscovery(t *testing.T) {
	cases := []struct {
		Name          string
		ExpectError   bool
		ServiceCruder func() mobile.ServiceCruder
		authChecker   func() mobile.AuthChecker
		client        mobile.ExternalHTTPRequester
		Validate      func(svs []*mobile.Service, t *testing.T)
	}{
		{
			Name: "test discover mobile services ok",
			client: &mock.Requester{
				Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader("")),
					}, nil
				},
			},
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
			svcs, err := is.DiscoverMobileServices(tc.ServiceCruder(), tc.authChecker(), tc.client)
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

func TestMobileServiceBindService(t *testing.T) {
	cases := []struct {
		Name              string
		SCClient          func() (mobile.SCCInterface, map[string]int)
		ExpectError       bool
		ServiceCruder     func() mobile.ServiceCruder
		TargetServiceID   string
		BindableServiceID string
	}{
		{
			Name:              "test bind ok",
			TargetServiceID:   mobile.ServiceNameSync,
			BindableServiceID: mobile.ServiceNameKeycloak,
			SCClient: func() (mobile.SCCInterface, map[string]int) {
				scc := mock.NewSCClient()
				calls := map[string]int{
					"BindToService": 1,
				}
				return scc, calls
			},
			ServiceCruder: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					name := action.(ktesting.GetAction).GetName()
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string][]byte{
							"uri":  []byte("http://" + name + ".com"),
							"name": []byte(name),
							"type": []byte(name),
						},
					}, nil
				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
		},
		{
			Name:              "test bind apiKeys ok",
			TargetServiceID:   mobile.ServiceNameSync,
			BindableServiceID: mobile.IntegrationAPIKeys,
			SCClient: func() (mobile.SCCInterface, map[string]int) {
				scc := mock.NewSCClient()
				calls := map[string]int{
					"AddMobileApiKeys": 1,
				}
				return scc, calls
			},
			ServiceCruder: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					name := action.(ktesting.GetAction).GetName()
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string][]byte{
							"uri":  []byte("http://" + name + ".com"),
							"name": []byte(name),
							"type": []byte(name),
						},
					}, nil
				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			is := integration.MobileService{}
			tscc, calls := tc.SCClient()
			err := is.BindMobileServices(tscc, tc.ServiceCruder(), tc.TargetServiceID, tc.BindableServiceID)
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one %s ", err)
			}
			for f, n := range calls {
				called := tscc.(*mock.SCClient).Called(f)
				if called != n {
					t.Fatalf("expected method %s to be called %v times but was %v times", f, n, called)
				}
			}
		})
	}
}

func TestMobileServiceUnBindService(t *testing.T) {
	cases := []struct {
		Name              string
		SCClient          func() (mobile.SCCInterface, map[string]int)
		ExpectError       bool
		ServiceCruder     func() mobile.ServiceCruder
		TargetServiceID   string
		BindableServiceID string
	}{
		{
			Name:              "test Unbind ok",
			TargetServiceID:   mobile.ServiceNameSync,
			BindableServiceID: mobile.ServiceNameKeycloak,
			SCClient: func() (mobile.SCCInterface, map[string]int) {
				scc := mock.NewSCClient()
				calls := map[string]int{
					"UnBindFromService": 1,
				}
				return scc, calls
			},
			ServiceCruder: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					name := action.(ktesting.GetAction).GetName()
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string][]byte{
							"uri":  []byte("http://" + name + ".com"),
							"name": []byte(name),
							"type": []byte(name),
						},
					}, nil
				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
		},
		{
			Name:              "test unbind apiKeys ok",
			TargetServiceID:   mobile.ServiceNameSync,
			BindableServiceID: mobile.IntegrationAPIKeys,
			SCClient: func() (mobile.SCCInterface, map[string]int) {
				scc := mock.NewSCClient()
				calls := map[string]int{
					"RemoveMobileApiKeys": 1,
				}
				return scc, calls
			},
			ServiceCruder: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					name := action.(ktesting.GetAction).GetName()
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string][]byte{
							"uri":  []byte("http://" + name + ".com"),
							"name": []byte(name),
							"type": []byte(name),
						},
					}, nil
				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			is := integration.MobileService{}
			tscc, calls := tc.SCClient()
			err := is.UnBindMobileServices(tscc, tc.ServiceCruder(), tc.TargetServiceID, tc.BindableServiceID)
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one %s ", err)
			}
			for f, n := range calls {
				called := tscc.(*mock.SCClient).Called(f)
				if called != n {
					t.Fatalf("expected method %s to be called %v times but was %v times", f, n, called)
				}
			}
		})
	}
}

func TestReadMobileServiceAndIntegrations(t *testing.T) {
	cases := []struct {
		Name          string
		ServiceCruder func() mobile.ServiceCruder
		AuthChecker   func() mobile.AuthChecker
		Requester     func() mobile.ExternalHTTPRequester
		ServiceName   string
		ExpectError   bool
		Validate      func(t *testing.T, ms *mobile.Service)
	}{
		{
			Name:        "test read mobile service and integrations ok",
			ServiceName: mobile.ServiceNameSync,
			AuthChecker: func() mobile.AuthChecker {
				return nil
			},
			Requester: func() mobile.ExternalHTTPRequester {
				return nil
			},
			ServiceCruder: func() mobile.ServiceCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					name := action.(ktesting.GetAction).GetName()
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string][]byte{
							"uri":  []byte("http://" + name + ".com"),
							"name": []byte(name),
							"type": []byte(name),
						},
					}, nil
				})
				client.AddReactor("list", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					list := v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Labels: map[string]string{"group": "mobileapp"},
								},
								Data: map[string][]byte{
									"uri":  []byte("http://fh-sync.com"),
									"name": []byte(mobile.IntegrationAPIKeys),
									"type": []byte(mobile.IntegrationAPIKeys),
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Labels: map[string]string{"group": "mobileapp"},
								},
								Data: map[string][]byte{
									"uri":  []byte("http://kc.com"),
									"name": []byte(mobile.ServiceNameKeycloak),
									"type": []byte(mobile.ServiceNameKeycloak),
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Labels: map[string]string{"group": "mobileapp"},
								},
								Data: map[string][]byte{
									"uri":  []byte("http://3scale.com"),
									"name": []byte(mobile.ServiceNameThreeScale),
									"type": []byte(mobile.ServiceNameThreeScale),
								},
							},
						},
					}
					return true, &list, nil

				})
				return data.NewMobileServiceRepo(client.CoreV1().Secrets("test"))
			},
			Validate: func(t *testing.T, ms *mobile.Service) {
				if nil == ms {
					t.Fatal("expected a mobile service but it was nil ")
				}
				if ms.Name != mobile.ServiceNameSync {
					t.Fatal("expected the sync service but got ", ms.Name)
				}
				if len(ms.Integrations) != 3 {
					t.Fatalf("Expected sync to have 3 integrations but got %v  ", len(ms.Integrations))
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			is := integration.MobileService{}
			ms, err := is.ReadMobileServiceAndIntegrations(tc.ServiceCruder(), tc.AuthChecker(), tc.ServiceName, tc.Requester())
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one %s ", err)
			}
			if nil != tc.Validate {
				tc.Validate(t, ms)
			}
		})
	}

}
