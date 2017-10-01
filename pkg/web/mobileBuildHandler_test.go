package web_test

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"bytes"
	"encoding/json"

	"io/ioutil"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/mock"
	"github.com/feedhenry/mcp-standalone/pkg/web"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	kfake "k8s.io/client-go/testing"
)

func setupMobileBuildHandler(kclient kubernetes.Interface, ocFake *kfake.Fake) http.Handler {
	r := web.NewRouter()
	logger := logrus.StandardLogger()
	if nil == kclient {
		kclient = &fake.Clientset{}
	}

	cb := &mock.ClientBuilder{
		Fakeclient: kclient,
	}
	ocClientBuilder := mock.NewOCClientBuilder("test", "test", "https://notthere.com", ocFake)
	repoBuilder := data.NewBuildsRepoBuilder(cb, ocClientBuilder, "test", "test")
	buildService := app.NewBuild()
	handler := web.NewBuildHandler(repoBuilder, buildService, logger)
	web.MobileBuildRoute(r, handler)
	return web.BuildHTTPHandler(r, nil)
}

func TestBuildHandlerCreate(t *testing.T) {
	cases := []struct {
		Name        string
		K8Client    func() kubernetes.Interface
		OCClient    func() *kfake.Fake
		ExpectError bool
		StatusCode  int
		MobileBuild *mobile.Build
		Validate    func(t *testing.T, ar *app.AppBuildCreatedResponse)
	}{
		{
			Name:       "test build create for private repo ok",
			StatusCode: 201,
			K8Client: func() kubernetes.Interface {
				c := &fake.Clientset{}
				c.AddReactor("create", "secrets", func(action kfake.Action) (handled bool, ret runtime.Object, err error) {
					obj := action.(kfake.CreateAction).GetObject()
					return true, obj, nil
				})
				return c
			},
			OCClient: func() *kfake.Fake {

				c := &kfake.Fake{}
				c.AddReactor("create", "buildconfig", func(action kfake.Action) (handled bool, ret runtime.Object, err error) {
					obj := action.(kfake.CreateAction).GetObject()
					return true, obj, nil
				})
				return c
			},
			MobileBuild: &mobile.Build{
				Name:  "mybuild",
				AppID: "myapp",
				GitRepo: &mobile.BuildGitRepo{
					Private: true,
					URI:     "git@git.com",
					Ref:     "master",
				},
			},
			Validate: func(t *testing.T, ar *app.AppBuildCreatedResponse) {
				if nil == ar {
					t.Fatal("expected a build creation response but got none")
				}
				if ar.PublicKey == "" {
					t.Fatal("expected a public key in the response but got none")
				}
				if ar.BuildID != "mybuild" {
					t.Fatalf("expected a build id to match : mybuild but got %s", ar.BuildID)
				}
			},
		},
		{
			Name:       "test build create for public repo ok",
			StatusCode: 201,
			K8Client: func() kubernetes.Interface {
				c := &fake.Clientset{}
				c.AddReactor("create", "secrets", func(action kfake.Action) (handled bool, ret runtime.Object, err error) {
					obj := action.(kfake.CreateAction).GetObject()
					return true, obj, nil
				})
				return c
			},
			OCClient: func() *kfake.Fake {

				c := &kfake.Fake{}
				c.AddReactor("create", "buildconfig", func(action kfake.Action) (handled bool, ret runtime.Object, err error) {
					obj := action.(kfake.CreateAction).GetObject()
					return true, obj, nil
				})
				return c
			},
			MobileBuild: &mobile.Build{
				Name:  "mybuild",
				AppID: "myapp",
				GitRepo: &mobile.BuildGitRepo{
					URI: "git@git.com",
					Ref: "master",
				},
			},
			Validate: func(t *testing.T, ar *app.AppBuildCreatedResponse) {
				if nil == ar {
					t.Fatal("expected a build creation response but got none")
				}
				if ar.PublicKey != "" {
					t.Fatal("did not expect a public key in the response but got one")
				}
				if ar.BuildID != "mybuild" {
					t.Fatalf("expected a build id to match : mybuild but got %s", ar.BuildID)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileBuildHandler(tc.K8Client(), tc.OCClient())
			server := httptest.NewServer(handler)
			defer server.Close()
			payload, err := json.Marshal(tc.MobileBuild)
			if err != nil {
				t.Fatal("failed to marshal json payload")
			}
			req, err := http.NewRequest("POST", server.URL+"/build", bytes.NewReader(payload))
			if err != nil {
				t.Fatalf("did not expect an error setting up request %s", err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("did not expect an error making request %s", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.StatusCode {
				t.Fatalf("expected status code %v but got %v ", tc.StatusCode, res.StatusCode)
			}
			if res.StatusCode == 201 {
				resBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("failed to read the response body %s", err)
				}
				response := &app.AppBuildCreatedResponse{}
				if err := json.Unmarshal(resBody, response); err != nil {
					t.Fatalf("did not expect an error unmarshalling the response body %s ", err)
				}
			}
		})
	}
}

func TestBuildHandlerGenerateKeys(t *testing.T) {
	cases := []struct {
		Name        string
		BuildID     string
		K8Client    func() kubernetes.Interface
		OCClient    func() *kfake.Fake
		ExpectError bool
		StatusCode  int
		Validate    func(res map[string]string, t *testing.T)
	}{
		{
			Name:    "test generate new keys ok",
			BuildID: "testbuild",
			K8Client: func() kubernetes.Interface {
				c := &fake.Clientset{}
				c.AddReactor("create", "secrets", func(action kfake.Action) (handled bool, ret runtime.Object, err error) {
					obj := action.(kfake.CreateAction).GetObject()
					return true, obj, nil
				})
				return c
			},
			OCClient: func() *kfake.Fake {
				c := &kfake.Fake{}
				return c
			},
			StatusCode: 201,
			Validate: func(res map[string]string, t *testing.T) {
				if res == nil {
					t.Fatal("expected a response body but got none")
				}
				if _, ok := res["name"]; !ok {
					t.Fatal("expected a name to be returned in the response")
				}
			},
		},
		{
			Name:    "test generate new keys fails if no buildID",
			BuildID: "",
			K8Client: func() kubernetes.Interface {
				c := &fake.Clientset{}
				return c
			},
			OCClient: func() *kfake.Fake {
				c := &kfake.Fake{}
				return c
			},
			StatusCode: 404,
		},
		{
			Name:    "test generate new keys fails if secret already exists",
			BuildID: "test",
			K8Client: func() kubernetes.Interface {
				c := &fake.Clientset{}
				c.AddReactor("create", "secrets", func(action kfake.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.NewConflict(schema.GroupResource{Resource: "", Group: ""}, "test", fmt.Errorf("this secret already exists "))
				})
				return c
			},
			OCClient: func() *kfake.Fake {
				c := &kfake.Fake{}
				return c
			},
			StatusCode: 409,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			handler := setupMobileBuildHandler(tc.K8Client(), tc.OCClient())
			server := httptest.NewServer(handler)
			defer server.Close()
			req, err := http.NewRequest("POST", server.URL+"/build/"+tc.BuildID+"/generatekeys", nil)
			if err != nil {
				t.Fatalf("did not expect an error setting up request %s", err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("did not expect an error making request %s", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.StatusCode {
				t.Fatalf("expected status code %v but got %v ", tc.StatusCode, res.StatusCode)
			}
			if res.StatusCode == 201 {
				resBod := map[string]string{}
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal("failed to read response body ", err)
				}
				if err := json.Unmarshal(data, &resBod); err != nil {
					t.Fatal("failed to unmarshal response body", err)
				}
				if tc.Validate != nil {
					tc.Validate(resBod, t)
				}
			}

		})
	}
}
