package testclient

import (
	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/testing"
)

var buildConfigsResource = schema.GroupVersionResource{Group: "", Version: "", Resource: "buildconfigs"}
var buildResource = schema.GroupVersionResource{Group: "", Version: "", Resource: "build"}
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

func NewFakeBuilds(ns string, fake *testing.Fake) *FakeBuilds {
	fb := &FakeBuilds{
		Namespace: ns,
		Fake:      fake,
	}
	if fb.Fake == nil {
		fb.Fake = &testing.Fake{}
	}
	return fb
}

type FakeBuilds struct {
	Fake      *testing.Fake
	Namespace string
}

func (c *FakeBuilds) Update(b *build.Build) (*build.Build, error) {
	rb, err := c.Fake.Invokes(testing.NewUpdateAction(buildResource, c.Namespace, b), b)
	if nil != rb {
		return rb.(*build.Build), err
	}
	return nil, err
}

func (c *FakeBuilds) Get(name string, options metav1.GetOptions) (*build.Build, error) {
	b, err := c.Fake.Invokes(testing.NewGetAction(buildResource, c.Namespace, name), &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{},
			Labels:      map[string]string{},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "build",
			APIVersion: "v1",
		},
	})
	if nil != b {
		return b.(*build.Build), err
	}
	return nil, err
}
