package data_test

import (
	"testing"

	"errors"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/testclient"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestBuildRepo_Create(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		BuildClient  func() client.BuildConfigInterface
		SecretClient func() corev1.SecretInterface
		Build        *mobile.Build
	}{
		{
			"test creating build succeeds",
			false,
			func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, nil
				})
				return fakeoc
			},
			func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				return kc.CoreV1().Secrets("test")
			},
			&mobile.Build{
				Name:  "test",
				AppID: "test",
				GitRepo: &mobile.BuildGitRepo{
					Private: true,
					URI:     "git@git.com",
					Ref:     "master",
				},
			},
		},
		{
			"test creating build fails when buildconfig not created",
			true,
			func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("failed to create")
				})
				return fakeoc
			},
			func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				return kc.CoreV1().Secrets("test")
			},
			&mobile.Build{
				Name:  "test",
				AppID: "test",
				GitRepo: &mobile.BuildGitRepo{
					Private: true,
					URI:     "git@git.com",
					Ref:     "master",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			buildRepo := data.NewBuildRepo(tc.BuildClient(), tc.SecretClient())
			err := buildRepo.Create(tc.Build)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}
		})
	}
}

func TestBuildRepo_AddBuildAsset(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		SecretClient func() corev1.SecretInterface
		BuildAsset   mobile.BuildAsset
	}{
		{
			Name: "test adding source credential asset ok",
			SecretClient: func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				kc.AddReactor("create", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					secret := action.(ktesting.CreateAction).GetObject().(*v1.Secret)
					if _, ok := secret.Labels["type"]; !ok || secret.Labels["type"] != string(mobile.BuildAssetTypeSourceCredential) {
						t.Fatalf("expected the sercret to have a lable with a type %s, %v", mobile.BuildAssetTypeSourceCredential, secret.Labels)
					}
					if _, ok := secret.Data["key"]; !ok {
						t.Fatalf("expected the secret to have data under with key:   key but it was missing %v", secret.Data)
					}
					return true, secret, nil
				})
				return kc.CoreV1().Secrets("test")
			},
			BuildAsset: mobile.BuildAsset{
				Name:    "myapp",
				AppName: "myapp",
				AssetData: map[string][]byte{
					"key": []byte("sdasdsadsadsa"),
				},
				Type:      mobile.BuildAssetTypeSourceCredential,
				BuildName: "mybuild",
			},
		},
		{
			Name: "test adding source credential fails when secret creation fails",
			SecretClient: func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				kc.AddReactor("create", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("failed to create secret")
				})
				return kc.CoreV1().Secrets("test")
			},
			BuildAsset: mobile.BuildAsset{
				Name:    "myapp",
				AppName: "myapp",
				AssetData: map[string][]byte{
					"key": []byte("sdasdsadsadsa"),
				},
				Type:      mobile.BuildAssetTypeSourceCredential,
				BuildName: "mybuild",
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			buildRepo := data.NewBuildRepo(nil, tc.SecretClient())
			secretName, err := buildRepo.AddBuildAsset(tc.BuildAsset)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}
			if !tc.ExpectError && secretName == "" {
				t.Fatal("expected a secret name to be returned but got none")
			}
		})
	}
}
