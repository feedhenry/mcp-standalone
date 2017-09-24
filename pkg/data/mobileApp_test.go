package data_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestReadMobileApp(t *testing.T) {
	cases := []struct {
		Name         string
		APIKeyClient func() corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
		ExpectError  bool
		Validate     func(app *mobile.App, t *testing.T)
	}{
		{
			Name: "test read mobile app ok",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.ConfigMap{
						Data: map[string]string{
							"name":       "app",
							"clientType": "android",
						},
					}, nil
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(a *mobile.App, t *testing.T) {
				if a == nil {
					t.Fatal("did not expect mobile app to be nil")
				}
				if a.ClientType != "android" {
					t.Fatalf("expected app type to be android but got %s", a.ClientType)
				}
			},
		},
		{
			Name: "test read mobile app fails when error",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("unexpected error")
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(a *mobile.App, t *testing.T) {
				if a != nil {
					t.Fatalf("expected mobile app to be nil but it wasn't %v", a)
				}
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(), data.DefaultMobileAppValidator{})
			app, err := appRepo.ReadByName("test")
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}
			tc.Validate(app, t)
		})
	}

}

func TestDeleteMobileApp(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		APIKeyClient func() corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
	}{
		{
			Name: "test delete mobile app ok",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("delete", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, nil
				})
				return c.CoreV1().ConfigMaps("test")
			},
		},
		{
			Name: "test delete mobile app fails on error",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("delete", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("failed to delete")
				})
				return c.CoreV1().ConfigMaps("test")
			},
			ExpectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(), data.DefaultMobileAppValidator{})
			err := appRepo.DeleteByName("test")
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one %v", err)
			}

		})
	}

}

func TestCreateMobileApp(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		APIKeyClient func() corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
		App          *mobile.App
	}{
		{
			Name: "test create mobile app ok",
			App: &mobile.App{
				Name:       "app",
				ClientType: "android",
				APIKey:     "akey",
				MetaData:   map[string]string{},
			},
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				return c.CoreV1().ConfigMaps("test")
			},
		},
		{
			Name: "test create mobile fails when invalid",
			App: &mobile.App{
				Name:       "app",
				ClientType: "nodroid",
			},
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				return c.CoreV1().ConfigMaps("test")
			},
			ExpectError: true,
		},
		{
			Name: "test create mobile fails when error returned from client",
			App: &mobile.App{
				Name:       "app",
				ClientType: "android",
				MetaData:   map[string]string{},
			},
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("create", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("already exists")
				})
				return c.CoreV1().ConfigMaps("test")
			},
			ExpectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(), data.DefaultMobileAppValidator{})
			err := appRepo.Create(tc.App)
			if tc.ExpectError && err == nil {
				t.Fatal("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an err but got one %v ", err)
			}
		})
	}
}

func TestListMobileApp(t *testing.T) {

	cases := []struct {
		Name         string
		ExpectError  bool
		APIKeyClient func() corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
		Validate     func(apps []*mobile.App, t *testing.T)
	}{
		{
			Name: "test list mobile apps ok",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("list", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					list := &v1.ConfigMapList{}
					list.Items = append(list.Items, v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobileapp"},
						},
						Data: map[string]string{
							"name":       "app",
							"clientType": "android",
						},
					},
						v1.ConfigMap{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{"group": "mobileapp"},
							},
							Data: map[string]string{
								"name":       "app2",
								"clientType": "iOS",
							},
						})
					return true, list, nil
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(apps []*mobile.App, t *testing.T) {
				if apps == nil {
					t.Fatal("expected apps but got nil")
				}
				if len(apps) != 2 {
					t.Fatalf("expected 2 apps but got %v ", len(apps))
				}
			},
		},
		{
			Name: "test list mobile apps doesn't list non mobile app",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("list", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					list := &v1.ConfigMapList{}
					list.Items = append(list.Items, v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
						},
						Data: map[string]string{
							"name": "something",
						},
					},
						v1.ConfigMap{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{"group": "notmobile"},
							},
							Data: map[string]string{
								"name": "something2",
							},
						})
					return true, list, nil
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(apps []*mobile.App, t *testing.T) {
				if apps == nil {
					t.Fatal("expected apps but got nil")
				}
				if len(apps) != 0 {
					t.Fatalf("expected 0 apps but got %v ", len(apps))
				}
			},
		},
		{
			Name:        "test list mobile fails on error",
			ExpectError: true,
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("list", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("failed to list")
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(apps []*mobile.App, t *testing.T) {
				if apps != nil {
					t.Fatalf("expected no apps but got %v", apps)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(), data.DefaultMobileAppValidator{})
			apps, err := appRepo.List()
			if tc.ExpectError && err == nil {
				t.Fatal("expexted an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an err but got one %v", err)
			}
			tc.Validate(apps, t)
		})
	}

}

func TestUpdateMobileApp(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		APIKeyClient func() corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
		App          *mobile.App
		Validate     func(app *mobile.App, t *testing.T)
	}{
		{
			App: &mobile.App{
				Name:       "app",
				ClientType: "iOS",
				Labels:     map[string]string{"group": "mobileapp"},
			},
			Name: "test update mobile apps clientType ok",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
						},
						Data: map[string]string{
							"name":       "app",
							"ClientTyoe": "iOS",
						},
					}, nil
				})
				c.AddReactor("update", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {

					return true, &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
						},
						Data: map[string]string{
							"name":       "app",
							"clientType": "android",
						},
					}, nil
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(a *mobile.App, t *testing.T) {
				if a == nil {
					t.Fatalf("expected an app but got nil ")
				}
				if a.Name != "app" {
					t.Fatalf("expected the app name to be the same but got %s ", a.Name)
				}
				if a.ClientType != "android" {
					t.Fatalf("expected the clientType to be android but got %s ", a.ClientType)
				}
			},
		},
		{
			App: &mobile.App{
				Name:       "app",
				ClientType: "not client",
				Labels:     map[string]string{"group": "mobileapp"},
			},
			Name: "test update mobile apps clientType fails when invalid",
			APIKeyClient: func() corev1.SecretInterface {
				c := fake.Clientset{}
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "mobile"},
						},
						Data: map[string]string{
							"name":       "app",
							"ClientTyoe": "iOS",
						},
					}, nil
				})
				return c.CoreV1().ConfigMaps("test")
			},
			Validate: func(a *mobile.App, t *testing.T) {
				if a != nil {
					t.Fatalf("expected no app but got one ")
				}
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(), data.DefaultMobileAppValidator{})
			app, err := appRepo.Update(tc.App)
			if tc.ExpectError && err == nil {
				t.Fatal("expexted an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an err but got one %v", err)
			}
			tc.Validate(app, t)
		})
	}
}

func TestUpdateAppAPIKeys(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		MobileApp    *mobile.App
		APIKeyClient func(t *testing.T) corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
	}{
		{
			Name: "test update api key ok",
			MobileApp: &mobile.App{
				ID:         "anID",
				APIKey:     "anAPIKey",
				Name:       "test",
				ClientType: "cordova",
				MetaData:   map[string]string{},
			},
			APIKeyClient: func(t *testing.T) corev1.SecretInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
						},
						Data: map[string][]byte{
							"client123": []byte("anAppKey"),
						},
					}, nil
				})
				c.AddReactor("update", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					sec := action.(ktesting.UpdateAction).GetObject().(*v1.Secret)
					key := sec.Data["anID"]
					if bytes.Compare(key, []byte("anAppKey")) == 0 {
						t.Fatal("Expected anID to have api key anAPIKey")
					}
					return true, nil, nil
				})
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				return c.CoreV1().ConfigMaps("test")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(t), data.DefaultMobileAppValidator{})
			err := appRepo.AddAPIKeyToMap(tc.MobileApp)
			if tc.ExpectError && err == nil {
				t.Fatal("expexted an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an err but got one %v", err)
			}
		})
	}
}

func TestRemoveAppAPIKeyByID(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		AppID        string
		APIKeyClient func(t *testing.T) corev1.SecretInterface
		Client       func() corev1.ConfigMapInterface
	}{
		{
			Name:  "test remove app api key ok",
			AppID: "anapp",
			APIKeyClient: func(t *testing.T) corev1.SecretInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"group": "notmobile"},
						},
						Data: map[string][]byte{
							"anapp": []byte("anAppKey"),
						},
					}, nil
				})
				c.AddReactor("update", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					sec := action.(ktesting.UpdateAction).GetObject().(*v1.Secret)
					if _, exists := sec.Data["anapp"]; exists {
						t.Fatal("Expected app anapp to be removed")
					}
					return true, nil, nil
				})
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				return c.CoreV1().ConfigMaps("test")
			},
		},
		{
			Name:  "test remove app api key error",
			AppID: "anapp",
			APIKeyClient: func(t *testing.T) corev1.SecretInterface {
				c := fake.Clientset{}
				c.AddReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("not found")
				})
				return c.CoreV1().Secrets("test")
			},
			Client: func() corev1.ConfigMapInterface {
				c := fake.Clientset{}
				return c.CoreV1().ConfigMaps("test")
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			appRepo := data.NewMobileAppRepo(tc.Client(), tc.APIKeyClient(t), data.DefaultMobileAppValidator{})
			err := appRepo.RemoveAPIKeyFromMap(tc.AppID)
			if tc.ExpectError && err == nil {
				t.Fatal("expexted an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an err but got one %v", err)
			}
		})
	}
}
