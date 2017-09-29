package app_test

import (
	"testing"

	"strings"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/testclient"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestBuildCreateAppBuild(t *testing.T) {
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
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					bc := action.(ktesting.CreateAction).GetObject().(*build.BuildConfig)
					if nil == bc {
						return true, nil, errors.New("no buildconfig passe")
					}
					if bc.Spec.Strategy.JenkinsPipelineStrategy.JenkinsfilePath != "Jenkinsfile" {
						return true, nil, errors.New("expected the JenkinsfilePath to be : Jenkinsfile but got " + bc.Spec.Strategy.JenkinsPipelineStrategy.JenkinsfilePath)
					}
					return true, bc, nil
				})
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
			Name: "test create app build with custom jenkinsfile location",
			Build: &mobile.Build{
				AppID: "test",
				Name:  "test",
				GitRepo: &mobile.BuildGitRepo{
					URI:             "git@git.com",
					Private:         true,
					JenkinsFilePath: "/build/Jenkinsfile",
				},
			},
			BuildClient: func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					bc := action.(ktesting.CreateAction).GetObject().(*build.BuildConfig)
					if nil == bc {
						return true, nil, errors.New("no buildconfig passe")
					}
					if bc.Spec.Strategy.JenkinsPipelineStrategy.JenkinsfilePath != "/build/Jenkinsfile" {
						return true, nil, errors.New("expected the JenkinsfilePath to be : /build/Jenkinsfile but got " + bc.Spec.Strategy.JenkinsPipelineStrategy.JenkinsfilePath)
					}
					return true, bc, nil
				})
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
			Name: "test create app build for public repo",
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
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					bc := action.(ktesting.CreateAction).GetObject().(*build.BuildConfig)
					if nil == bc {
						return true, nil, errors.New("no buildconfig passe")
					}
					if bc.Spec.Strategy.JenkinsPipelineStrategy.JenkinsfilePath != "Jenkinsfile" {
						return true, nil, errors.New("expected the JenkinsfilePath to be : Jenkinsfile but got " + bc.Spec.Strategy.JenkinsPipelineStrategy.JenkinsfilePath)
					}
					return true, bc, nil
				})
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
		{
			Name: "test create app build fails when cant create buildconfig",
			Build: &mobile.Build{
				AppID: "test",
				Name:  "test",
				GitRepo: &mobile.BuildGitRepo{
					URI:     "git@git.com",
					Private: false,
				},
			},
			ExpectError: true,
			BuildClient: func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("faield to create buildconfig")
				})
				return fakeoc
			},
			SecretClient: func() corev1.SecretInterface {
				fakes := &fake.Clientset{}
				return fakes.CoreV1().Secrets("test")
			},
			Validate: func(t *testing.T, br *app.AppBuildCreatedResponse) {
				if nil != br {
					t.Fatalf("expected no AppBuildCreatedResponse but got %v ", br)
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

func TestBuildCreateBuildSrcKeySecret(t *testing.T) {
	cases := []struct {
		Name         string
		ExpectError  bool
		BuildClient  func() client.BuildConfigInterface
		SecretClient func() corev1.SecretInterface
	}{
		{
			Name: "test creating src key secret ok",
			BuildClient: func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				return fakeoc
			},
			SecretClient: func() corev1.SecretInterface {
				fakec := &fake.Clientset{}
				fakec.AddReactor("create", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					secret := action.(ktesting.CreateAction).GetObject().(*v1.Secret)
					return true, secret, nil

				})
				return fakec.CoreV1().Secrets("test")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {

			br := data.NewBuildRepo(tc.BuildClient(), tc.SecretClient())
			buildService := app.NewBuild()
			secretName, pubKey, err := buildService.CreateBuildSrcKeySecret(br, "test", "test")
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an err but got none!")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect and err but got one %v", err)
			}
			if !tc.ExpectError && secretName == "" {
				t.Fatal("expected a secretname to be returned but found none")
			}

			if !tc.ExpectError && string(pubKey) == "" {
				t.Fatalf("expected a public key to be returned but got none")
			}

		})
	}
}
