package data_test

import (
	"testing"

	"errors"

	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/testclient"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestBuildRepo_Create(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectError       bool
		BuildConfigClient func() client.BuildConfigInterface
		BuildClient       func() client.BuildInterface
		SecretClient      func() corev1.SecretInterface
		Build             *mobile.BuildConfig
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
			func() client.BuildInterface {
				fakeoc := testclient.NewFakeBuilds("test", nil)
				return fakeoc

			},
			func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				return kc.CoreV1().Secrets("test")
			},
			&mobile.BuildConfig{
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
			func() client.BuildInterface {
				fakeoc := testclient.NewFakeBuilds("test", nil)
				return fakeoc

			},
			func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				return kc.CoreV1().Secrets("test")
			},
			&mobile.BuildConfig{
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
			buildRepo := data.NewBuildRepo(tc.BuildConfigClient(), tc.BuildClient(), tc.SecretClient())
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

func TestBuildRepo(t *testing.T) {
	cases := []struct {
		Name              string
		ExpectError       bool
		BuildName         string
		BuildConfigClient func() client.BuildConfigInterface
		BuildClient       func() client.BuildInterface
		SecretClient      func() corev1.SecretInterface
	}{
		{
			Name:      "test building mobile app",
			BuildName: "buildname",
			BuildConfigClient: func() client.BuildConfigInterface {
				fakeoc := testclient.NewFakeBuildConfigs("test", nil)
				fakeoc.Fake.AddReactor("create", "buildconfigs", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, nil
				})
				return fakeoc
			},
			BuildClient: func() client.BuildInterface {
				fakeoc := testclient.NewFakeBuilds("test", nil)
				return fakeoc

			},
			SecretClient: func() corev1.SecretInterface {
				kc := &fake.Clientset{}
				return kc.CoreV1().Secrets("test")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			buildRepo := data.NewBuildRepo(tc.BuildConfigClient(), tc.BuildClient(), tc.SecretClient())
			err := buildRepo.BuildApp(tc.Name)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s", err)
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
			buildRepo := data.NewBuildRepo(nil, nil, tc.SecretClient())
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

func TestBuildRepoAddDownload(t *testing.T) {
	cases := []struct {
		Name        string
		ExpectError bool
		BuildClient func(t *testing.T) client.BuildInterface
		Download    *mobile.BuildDownload
		BuildName   string
	}{
		{
			Name: "test adding a build download is ok ",
			BuildClient: func(t *testing.T) client.BuildInterface {
				fakeoc := testclient.NewFakeBuilds("test", nil)
				fakeoc.Fake.AddReactor("update", "build", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					b := action.(ktesting.UpdateAction).GetObject().(*build.Build)
					if nil == b {
						t.Fatal("expected a build but it was nil")
					}
					if _, ok := b.Annotations["downloadURL"]; !ok {
						t.Fatal("expected a downloadURL but it wasnt present")
					}
					if _, ok := b.Annotations["downloadExpires"]; ok {
						t.Fatal("expected an expires but there was none")
					}

					return true, b, nil
				})
				return fakeoc
			},
			Download: &mobile.BuildDownload{
				URL:     "https://mcp.com/buid/test?token=asdadasd",
				Expires: 12334343434,
			},
			BuildName: "test",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			buildRepo := data.NewBuildRepo(nil, tc.BuildClient(t), nil)
			err := buildRepo.AddDownload(tc.BuildName, tc.Download)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}
		})
	}
}

func TestBuildRepoStatus(t *testing.T) {
	cases := []struct {
		Name        string
		BuildClient func() client.BuildInterface
		ExpectError bool
		Validate    func(status *mobile.BuildStatus, err error, t *testing.T)
	}{
		{
			Name: "test getting build status ok",
			BuildClient: func() client.BuildInterface {
				fakeoc := testclient.NewFakeBuilds("test", nil)
				fakeoc.Fake.AddReactor("get", "build", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					b := &build.Build{
						ObjectMeta: meta_v1.ObjectMeta{
							Annotations: map[string]string{
								"openshift.io/jenkins-status-json": `{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/wfapi/describe"},"artifacts":{"href":"//job/localmcp-androidapp1/1/wfapi/artifacts"}},"id":"1","name":"#1","status":"SUCCESS","startTimeMillis":1506071344983,"endTimeMillis":0,"durationMillis":0,"queueDurationMillis":157994,"pauseDurationMillis":0,"stages":[{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/6/wfapi/describe"}},"id":"6","name":"Checkout","execNode":"","status":"SUCCESS","startTimeMillis":1506071502977,"durationMillis":7089,"pauseDurationMillis":0,"stageFlowNodes":[{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/7/wfapi/describe"}},"id":"7","name":"General SCM","execNode":"","status":"SUCCESS","startTimeMillis":1506071503267,"durationMillis":6792,"pauseDurationMillis":0,"parentNodes":["6"]}]},{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/11/wfapi/describe"}},"id":"11","name":"Prepare","execNode":"","status":"SUCCESS","startTimeMillis":1506071510078,"durationMillis":81130,"pauseDurationMillis":0,"stageFlowNodes":[{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/12/wfapi/describe"}},"id":"12","name":"Shell Script","execNode":"","status":"SUCCESS","startTimeMillis":1506071512564,"durationMillis":78639,"pauseDurationMillis":0,"parentNodes":["11"]}]},{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/16/wfapi/describe"}},"id":"16","name":"Build","execNode":"","status":"SUCCESS","startTimeMillis":1506071591273,"durationMillis":175102,"pauseDurationMillis":0,"stageFlowNodes":[{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/17/wfapi/describe"}},"id":"17","name":"Shell Script","execNode":"","status":"SUCCESS","startTimeMillis":1506071591384,"durationMillis":174986,"pauseDurationMillis":0,"parentNodes":["16"]}]},{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/21/wfapi/describe"}},"id":"21","name":"Archive","execNode":"","status":"SUCCESS","startTimeMillis":1506071766385,"durationMillis":2588,"pauseDurationMillis":0,"stageFlowNodes":[{"_links":{"self":{"href":"http://jenkins-localmcp.192.168.37.1.nip.io/job/localmcp-androidapp1/1/execution/node/22/wfapi/describe"}},"id":"22","name":"General Build Step","execNode":"","status":"SUCCESS","startTimeMillis":1506071766665,"durationMillis":2305,"pauseDurationMillis":0,"parentNodes":["21"]}]}]}`,
							},
							Name: "test",
						},
						Status: build.BuildStatus{
							Phase: build.BuildPhaseComplete,
						},
					}
					return true, b, nil
				})
				return fakeoc
			},
			Validate: func(status *mobile.BuildStatus, err error, t *testing.T) {
				if nil == status {
					t.Fatal("expected a build status but got none")
				}
				if status.Phase != "Complete" {
					t.Fatalf("expected phase to be Complete but got: %s", status.Phase)
				}
				if status.Links.Artifacts.Href == "" {
					t.Fatal("expected an artifact href but got none")
				}
			},
		},
		{
			Name: "test getting build fails with not found when no status",
			BuildClient: func() client.BuildInterface {
				fakeoc := testclient.NewFakeBuilds("test", nil)
				fakeoc.Fake.AddReactor("get", "build", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					b := &build.Build{
						ObjectMeta: meta_v1.ObjectMeta{
							Annotations: map[string]string{},
							Name:        "test",
						},
						Status: build.BuildStatus{
							Phase: build.BuildPhaseComplete,
						},
					}
					return true, b, nil
				})
				return fakeoc
			},
			ExpectError: true,
			Validate: func(status *mobile.BuildStatus, err error, t *testing.T) {
				if nil != status {
					t.Fatal("did not expect a build status but got one")
				}
				if !data.IsNotFoundErr(err) {
					t.Fatal("expected a not found error ")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			buildRepo := data.NewBuildRepo(nil, tc.BuildClient(), nil)
			bs, err := buildRepo.Status("test")
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}
			if nil != tc.Validate {
				tc.Validate(bs, err, t)
			}
		})
	}

}
