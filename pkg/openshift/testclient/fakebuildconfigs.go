package testclient

import (
	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/testing"
)

var buildConfigsResource = schema.GroupVersionResource{Group: "", Version: "", Resource: "buildconfigs"}
var buildConfigsKind = schema.GroupVersionKind{Group: "", Version: "", Kind: "BuildConfig"}

// FakeBuildConfigs implements BuildConfigInterface. Meant to be embedded into a struct to get a default
// implementation. This makes faking out just the methods you want to test easier.
type FakeBuildConfigs struct {
	Fake      *testing.Fake
	Namespace string
}

func NewFakeBuildConfigs(ns string, fake *testing.Fake) *FakeBuildConfigs {
	fb := &FakeBuildConfigs{
		Namespace: ns,
		Fake:      fake,
	}
	if fb.Fake == nil {
		fb.Fake = &testing.Fake{}
	}
	return fb
}

func (c *FakeBuildConfigs) Create(config *build.BuildConfig) (*build.BuildConfig, error) {
	bc, err := c.Fake.Invokes(testing.NewCreateAction(buildConfigsResource, c.Namespace, config), config)
	if nil != bc {
		return bc.(*build.BuildConfig), err
	}
	return nil, err
}

func (c *FakeBuildConfigs) Update(config *build.BuildConfig) (*build.BuildConfig, error) {
	bc, err := c.Fake.Invokes(testing.NewUpdateAction(buildConfigsResource, c.Namespace, config), config)
	if nil != bc {
		return bc.(*build.BuildConfig), err
	}
	return nil, err
}
