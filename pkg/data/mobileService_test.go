package data_test

import (
	"errors"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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
									"name": []byte("fh-sync"),
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
				if s.Name != "fh-sync" {
					t.Fatalf("expected the service to be fh-sync but got %s ", s.Name)
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
				return a.GetName() == "fh-sync"
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
