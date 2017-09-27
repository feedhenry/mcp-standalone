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
