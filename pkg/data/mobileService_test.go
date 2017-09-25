package data_test

import (
	"errors"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mock"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestListMobileServices(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() kubernetes.Interface
		ExpectError bool
		Validate    func(svs []*mobile.Service, t *testing.T)
	}{
		{
			Name: "test list mobile services ok",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(a ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"name": []byte("fh-sync-server"),
								},
							},
							{
								Data: map[string][]byte{
									"now":        []byte("something"),
									"completely": []byte("different"),
								},
							},
						},
					}, nil
				})
				return client
			},
			Validate: func(svs []*mobile.Service, t *testing.T) {
				if nil == svs {
					t.Fatal("expected services but got none")
				}
				if len(svs) != 1 {
					t.Fatalf("expected 1 service to be returned but got %v ", len(svs))
				}
				s := svs[0]
				if s.Name != "fh-sync-server" {
					t.Fatalf("expected the service to be fh-sync-server but got %s ", s.Name)
				}
			},
		},
		{
			Name: "test list mobile services fails on error",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(a ktesting.Action) (bool, runtime.Object, error) {
					return true, nil, errors.New("fatal error")
				})
				return client
			},
			ExpectError: true,
			Validate: func(svs []*mobile.Service, t *testing.T) {
				if nil != svs {
					t.Fatal("expected no services but got some", svs)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client().CoreV1().Secrets("test")
			mobileRepo := data.NewMobileServiceRepo(client)
			svc, err := mobileRepo.List(func(a mobile.Attributer) bool {
				return a.GetName() == "fh-sync-server"
			})
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %v ", err)
			}
			tc.Validate(svc, t)
		})
	}
}

func TestListMobileServiceConfigs(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() kubernetes.Interface
		ExpectError bool
		Validate    func(svs []*mobile.ServiceConfig, t *testing.T)
	}{
		{
			Name: "test listing service configs is ok",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(a ktesting.Action) (bool, runtime.Object, error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								Data: map[string][]byte{
									"uri":  []byte("http://test.com"),
									"name": []byte("fh-sync-server"),
								},
							},
							mock.KeycloakSecret(),
						},
					}, nil
				})
				return client
			},
			Validate: func(svc []*mobile.ServiceConfig, t *testing.T) {
				if len(svc) != 2 {
					t.Fatalf("expected 2 service configs but got %v ", len(svc))
				}
				foundSync := false
				foundKeyCloak := false
				for _, sc := range svc {
					if sc.Name == "fh-sync-server" {
						foundSync = true
					}
					if sc.Name == "keycloak" {
						foundKeyCloak = true
					}
				}
				if !foundSync || !foundKeyCloak {
					t.Fatal("expected to find keycloak and sync configs but didn't ", svc)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := tc.Client().CoreV1().Secrets("test")
			mobileRepo := data.NewMobileServiceRepo(client)
			svc, err := mobileRepo.ListConfigs(func(a mobile.Attributer) bool {
				return a.GetName() == "fh-sync-server" || a.GetName() == "keycloak"
			})
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %v ", err)
			}
			tc.Validate(svc, t)
		})
	}
}

func TestReadMobileService(t *testing.T) {
	cases := []struct {
		Name            string
		Client          func() corev1.SecretInterface
		ExpectError     bool
		ExpectedService *mobile.Service
		Validate        func(expectd *mobile.Service, actual *mobile.Service, t *testing.T)
	}{
		{
			Name: "should read mobile service ok",
			ExpectedService: &mobile.Service{
				ID:       "somesecretID",
				Name:     "fh-sync-server",
				Host:     "https://test.com",
				External: true,
			},
			Client: func() corev1.SecretInterface {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
							Name:   "somesecretID",
						},
						Data: map[string][]byte{
							"name": []byte("fh-sync-server"),
							"uri":  []byte("https://test.com"),
						},
					}, nil
				})
				return client.CoreV1().Secrets("test")
			},
			Validate: func(expected, svc *mobile.Service, t *testing.T) {
				if nil == svc {
					t.Fatal("expected a mobile service but it was nil ")
				}
				if expected.Name != svc.Name {
					t.Fatalf("expected the mobile service name to be %s but got %s ", expected.Name, svc.Name)
				}
				if expected.Host != svc.Host {
					t.Fatalf("expected the mobile service host to be %s but got %s ", expected.Host, svc.Host)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mobileRepo := data.NewMobileServiceRepo(tc.Client())
			ms, err := mobileRepo.Read(tc.ExpectedService.ID)
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %v ", err)
			}
			tc.Validate(tc.ExpectedService, ms, t)
		})
	}
}

func TestMobileServiceDisplayName(t *testing.T) {
	cases := []struct {
		Name                string
		Client              func() corev1.SecretInterface
		ExpectedDisplayName string
	}{
		{
			Name: "should use name as display name ok",
			Client: func() corev1.SecretInterface {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
							Name:   "aService",
						},
						Data: map[string][]byte{
							"name": []byte("fh-sync-server"),
							"uri":  []byte("https://test.com"),
						},
					}, nil
				})
				return client.CoreV1().Secrets("test")
			},
			ExpectedDisplayName: "fh-sync-server",
		},
		{
			Name: "should use displayName as display name ok",
			Client: func() corev1.SecretInterface {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
							Name:   "aService",
						},
						Data: map[string][]byte{
							"name":        []byte("fh-sync-server"),
							"displayName": []byte("Sync Server"),
							"uri":         []byte("https://test.com"),
						},
					}, nil
				})
				return client.CoreV1().Secrets("test")
			},
			ExpectedDisplayName: "Sync Server",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mobileRepo := data.NewMobileServiceRepo(tc.Client())
			ms, _ := mobileRepo.Read("test")
			if tc.ExpectedDisplayName != ms.DisplayName {
				t.Fatalf("expected the display name to be %s but got %s ", tc.ExpectedDisplayName, ms.DisplayName)
			}
		})
	}
}

func TestMobileServiceRepo_Delete(t *testing.T) {
	cases := []struct {
		Name        string
		ExpectError bool
		Client      func() corev1.SecretInterface
		ServiceName string
	}{
		{
			Name:        "test delete service ok",
			ServiceName: "myservice",
			Client: func() corev1.SecretInterface {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
							Name:   "myservice",
						},
						Data: map[string][]byte{
							"name": []byte("myservice"),
							"uri":  []byte("https://test.com"),
						},
					}, nil
				})
				return client.CoreV1().Secrets("test")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mobileRepo := data.NewMobileServiceRepo(tc.Client())
			err := mobileRepo.Delete(tc.ServiceName)
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %v ", err)
			}
		})
	}
}
