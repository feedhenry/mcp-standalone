package app_test

import (
	"errors"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestCreateApp(t *testing.T) {
	cases := []struct {
		Name          string
		ExpectError   bool
		ServiceCruder func() mobile.ServiceCruder
		AppCruder     func() mobile.AppCruder
		MobileApp     *mobile.App
	}{
		{
			Name: "test create app ok",
			AppCruder: func() mobile.AppCruder {
				client := &fake.Clientset{}
				return data.NewMobileAppRepo(client.CoreV1().ConfigMaps("test"), client.CoreV1().Secrets("test"), nil)
			},
			MobileApp: &mobile.App{
				Name:       "test",
				ClientType: "cordova",
				MetaData:   map[string]string{},
			},
		},
		{
			Name:        "test create app error",
			ExpectError: true,
			AppCruder: func() mobile.AppCruder {
				client := &fake.Clientset{}
				return data.NewMobileAppRepo(client.CoreV1().ConfigMaps("test"), client.CoreV1().Secrets("test"), nil)
			},
			MobileApp: &mobile.App{
				Name:     "test",
				MetaData: map[string]string{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			service := &app.Service{}
			err := service.Create(tc.AppCruder(), tc.MobileApp)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none!")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect and err but got one %v", err)
			}
		})
	}
}

func TestRemoveAppByID(t *testing.T) {
	cases := []struct {
		Name        string
		ExpectError bool
		AppCruder   func() mobile.AppCruder
		AppID       string
	}{
		{
			Name: "test remove app ok",
			AppCruder: func() mobile.AppCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
						},
						Data: map[string][]byte{
							"anapp": []byte("anAppKey"),
						},
					}, nil
				})
				return data.NewMobileAppRepo(client.CoreV1().ConfigMaps("test"), client.CoreV1().Secrets("test"), nil)
			},
			AppID: "anapp",
		},
		{
			Name:        "test remove app error",
			ExpectError: true,
			AppCruder: func() mobile.AppCruder {
				client := &fake.Clientset{}
				client.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("secret doesn't exist")
				})
				return data.NewMobileAppRepo(client.CoreV1().ConfigMaps("test"), client.CoreV1().Secrets("test"), nil)
			},
			AppID: "anotherapp",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			service := &app.Service{}
			err := service.Delete(tc.AppCruder(), tc.AppID)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none!")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect and err but got one %v", err)
			}
		})
	}
}
