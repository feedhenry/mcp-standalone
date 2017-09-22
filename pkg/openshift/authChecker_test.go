package openshift

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	testHost      = "test-host"
	testUser      = &User{User: "Phil", Groups: []string{"group1", "group2"}}
	testToken     = "test-token"
	testUserRepo  = newMockUserRepo(testUser, nil)
	testResource  = "deployments-test"
	testNamespace = "namespace-test"
)

type responseBody struct {
	Kind            string   `json:"kind"`
	ApiVersion      string   `json:"apiVersion"`
	Namespace       string   `json:"namespace"`
	Users           []string `json:"users"`
	Groups          []string `json:"groups"`
	EvaluationError string   `json:"evaluationError"`
}

func newMockResponse(u, g []string, namespace, evalError string) *responseBody {
	users := []string{
		"system:admin",
		"system:serviceaccount:ansible-service-broker:asb",
		"system:serviceaccount:default:pvinstaller",
		"system:serviceaccount:kube-system:generic-garbage-collector",
		"system:serviceaccount:openshift-infra:image-trigger-controller",
		"system:serviceaccount:openshift-infra:template-instance-controller",
		"system:serviceaccount:openshift-infra:unidling-controller",
	}
	users = append(users, u...)
	groups := []string{
		"system:cluster-admins",
		"system:masters",
	}
	groups = append(groups, g...)

	return &responseBody{
		Kind:            "ResourceAccessReviewResponse",
		ApiVersion:      "v1",
		Users:           users,
		Groups:          groups,
		Namespace:       namespace,
		EvaluationError: evalError,
	}
}

type mockExternalHTTPRequester struct {
	ResponseBody  io.ReadCloser
	ResponseError error
	ResponseCode  int
}

func (mehr *mockExternalHTTPRequester) Do(req *http.Request) (*http.Response, error) {
	res := &http.Response{Body: mehr.ResponseBody, StatusCode: mehr.ResponseCode}
	return res, mehr.ResponseError
}

func (mehr *mockExternalHTTPRequester) Get(url string) (*http.Response, error) {
	res := &http.Response{Body: mehr.ResponseBody, StatusCode: mehr.ResponseCode}
	return res, mehr.ResponseError
}

func newMockExternalHTTPRequester(body io.ReadCloser, status int, err error) mobile.ExternalHTTPRequester {
	return &mockExternalHTTPRequester{ResponseBody: body, ResponseError: err, ResponseCode: status}
}

type mockUserRepo struct {
	User mobile.User
	Err  error
}

func (mur *mockUserRepo) GetUser() (mobile.User, error) {
	return mur.User, mur.Err
}

func newMockUserRepo(user mobile.User, err error) *mockUserRepo {
	return &mockUserRepo{User: user, Err: err}
}

func TestNewAuthCheckerBuilder(t *testing.T) {
	acb := NewAuthCheckerBuilder(testHost)
	if reflect.TypeOf(acb).String() != "*openshift.AuthCheckerBuilder" {
		t.Fatalf("expected '*openshift.AuthCheckerBuilder' got '%s'", reflect.TypeOf(acb).String())
	}
}

func TestAuthCheckerBuilder_Build(t *testing.T) {
	ac := NewAuthCheckerBuilder(testHost).Build()
	if reflect.TypeOf(ac).String() != "*openshift.AuthChecker" {
		t.Fatalf("exected '*openshift.AuthChecker' but got '%s'", reflect.TypeOf(ac).String())
	}
	ac = NewAuthCheckerBuilder(testHost).WithUserRepo(testUserRepo).Build()
	if reflect.TypeOf(ac).String() != "*openshift.AuthChecker" {
		t.Fatalf("exected '*openshift.AuthChecker' but got '%s'", reflect.TypeOf(ac).String())
	}
	ac = NewAuthCheckerBuilder(testHost).WithToken(testToken).Build()
	if reflect.TypeOf(ac).String() != "*openshift.AuthChecker" {
		t.Fatalf("exected '*openshift.AuthChecker' but got '%s'", reflect.TypeOf(ac).String())
	}
	ac = NewAuthCheckerBuilder(testHost).WithToken(testToken).WithUserRepo(testUserRepo).Build()
	if reflect.TypeOf(ac).String() != "*openshift.AuthChecker" {
		t.Fatalf("exected '*openshift.AuthChecker' but got '%s'", reflect.TypeOf(ac).String())
	}
	ac = NewAuthCheckerBuilder(testHost).IgnoreCerts().Build()
	if reflect.TypeOf(ac).String() != "*openshift.AuthChecker" {
		t.Fatalf("exected '*openshift.AuthChecker' but got '%s'", reflect.TypeOf(ac).String())
	}
}

func TestAuthChecker_Check(t *testing.T) {
	testcases := []struct {
		TestName   string
		UserRepo   mobile.UserRepo
		Users      []string
		Groups     []string
		Namespace  string
		EvalError  string
		StatusCode int
		ResError   error
		Validation func(res bool, err error, t *testing.T) error
	}{
		{
			TestName:   "User is authorized",
			UserRepo:   testUserRepo,
			Users:      []string{testUser.User},
			Groups:     testUser.Groups,
			Namespace:  testNamespace,
			StatusCode: http.StatusCreated,
			ResError:   nil,
			EvalError:  "",
			Validation: func(res bool, err error, t *testing.T) error {
				if err != nil {
					return errors.New(fmt.Sprintf("unexpected error: '%+v'", err))
				}
				if !res {
					return errors.New(fmt.Sprintf("got false, expected true"))
				}
				return nil
			},
		},
		{
			TestName:   "User is not authorized to edit deployments",
			UserRepo:   testUserRepo,
			Users:      []string{},
			Groups:     []string{},
			Namespace:  testNamespace,
			StatusCode: http.StatusUnauthorized,
			ResError:   nil,
			EvalError:  "",
			Validation: func(res bool, err error, t *testing.T) error {
				if err == nil {
					return errors.New(fmt.Sprintf("expected error: 'access denied (401)' but got no error"))
				}
				if res {
					return errors.New(fmt.Sprintf("got true, expected false"))
				}
				return nil
			},
		},
		{
			TestName:   "User is not authorized to create permissions check",
			UserRepo:   testUserRepo,
			Users:      []string{},
			Groups:     []string{},
			Namespace:  testNamespace,
			StatusCode: http.StatusForbidden,
			ResError:   nil,
			EvalError:  "",
			Validation: func(res bool, err error, t *testing.T) error {
				if err != nil {
					return errors.New(fmt.Sprintf("unexpected error: '%+v'", err))
				}
				if res {
					return errors.New(fmt.Sprintf("got true, expected false"))
				}
				return nil
			},
		},
		{
			TestName:   "User is in authorized group",
			UserRepo:   testUserRepo,
			Users:      []string{},
			Groups:     testUser.Groups,
			Namespace:  testNamespace,
			StatusCode: http.StatusCreated,
			ResError:   nil,
			EvalError:  "",
			Validation: func(res bool, err error, t *testing.T) error {
				if err != nil {
					return errors.New(fmt.Sprintf("unexpected error: '%+v'", err))
				}
				if !res {
					return errors.New(fmt.Sprintf("got false, expected true"))
				}
				return nil
			},
		},
		{
			TestName:   "User not in authorized users or groups",
			UserRepo:   testUserRepo,
			Users:      []string{},
			Groups:     []string{},
			Namespace:  testNamespace,
			StatusCode: http.StatusCreated,
			ResError:   nil,
			EvalError:  "",
			Validation: func(res bool, err error, t *testing.T) error {
				if err != nil {
					return errors.New(fmt.Sprintf("unexpected error: '%+v'", err))
				}
				if res {
					return errors.New(fmt.Sprintf("got true, expected false"))
				}
				return nil
			},
		},
	}

	for _, tc := range testcases {
		ac := NewAuthCheckerBuilder(testHost).WithToken(testToken).WithUserRepo(tc.UserRepo).Build()
		resBody := newMockResponse(tc.Users, tc.Groups, tc.Namespace, tc.EvalError)
		resBytes, err := json.Marshal(resBody)
		if err != nil {
			t.Fatalf("unexpected error creating mock response: '%+v'", err)
		}
		res, err := ac.Check(testResource, testNamespace, newMockExternalHTTPRequester(ioutil.NopCloser(bytes.NewReader(resBytes)), tc.StatusCode, tc.ResError))
		err = tc.Validation(res, err, t)
		if err != nil {
			t.Fatalf("error in testcase '%s': '%+v'", tc.TestName, err)
		}
	}
}
