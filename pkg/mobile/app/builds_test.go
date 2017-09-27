package app_test

import (
	"testing"

	"strings"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/testclient"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func TestBuild_CreateAppBuild(t *testing.T) {
	cases := []struct {
		Name         string
		Build        *mobile.Build
		ExpectError  bool
		BuildClient  func() client.BuildConfigInterface
		SecretClient func() corev1.SecretInterface
		Validate     func(t *testing.T, br *app.AppBuildCreatedResponse)
	}{
		{
			Name: "test create app build with new public private key",
			Build: &mobile.Build{
				AppID: "test",
				Name:  "test",
				GitRepo: &mobile.BuildGitRepo{
					URI:     "git@git.com",
					Private: true,
				},
			},
			BuildClient: func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				return fakeoc
			},
			SecretClient: func() corev1.SecretInterface {
				fakec := &fake.Clientset{}
				return fakec.CoreV1().Secrets("test")
			},
			Validate: func(t *testing.T, br *app.AppBuildCreatedResponse) {
				if nil == br {
					t.Fatal("expected an app build response but got nil")
				}
				if br.PublicKey == "" {
					t.Fatal("expected a public key but it was empty")
				}
				if !strings.Contains(br.PublicKey, "PUBLIC KEY") {
					t.Fatal("expected a public key but did not find public key comment ", br.PublicKey)
				}
			},
		},
		{
			Name: "test create app build with no public private key",
			Build: &mobile.Build{
				AppID: "test",
				Name:  "test",
				GitRepo: &mobile.BuildGitRepo{
					URI:     "git@git.com",
					Private: false,
				},
			},
			BuildClient: func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				return fakeoc
			},
			SecretClient: func() corev1.SecretInterface {
				fakes := &fake.Clientset{}
				return fakes.CoreV1().Secrets("test")
			},
			Validate: func(t *testing.T, br *app.AppBuildCreatedResponse) {
				if nil == br {
					t.Fatal("expected an app build response but got nil")
				}
				if br.PublicKey != "" {
					t.Fatalf("did not expect a public key but found %s ", br.PublicKey)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			bc := data.NewBuildRepo(tc.BuildClient(), tc.SecretClient())
			buildService := app.NewBuild()
			br, err := buildService.CreateAppBuild(bc, tc.Build)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none!")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect and err but got one %v", err)
			}
			if tc.Validate != nil {
				tc.Validate(t, br)
			}
		})
	}
}
