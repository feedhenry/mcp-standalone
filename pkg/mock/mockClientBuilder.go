package mock

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/testclient"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/testing"
)

// ClientBuilder is a mock kubernetes client builder
type ClientBuilder struct {
	token      string
	namespace  string
	host       string
	Fakeclient kubernetes.Interface
}

func (cb *ClientBuilder) WithToken(token string) mobile.K8ClientBuilder {
	return &ClientBuilder{token: token, host: cb.host, namespace: cb.namespace, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) WithNamespace(ns string) mobile.K8ClientBuilder {
	return &ClientBuilder{token: cb.token, host: cb.host, namespace: ns, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) WithHost(host string) mobile.K8ClientBuilder {
	return &ClientBuilder{token: cb.token, host: host, namespace: cb.namespace, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) WithHostAndNamespace(host, ns string) mobile.K8ClientBuilder {
	return &ClientBuilder{token: cb.token, host: host, namespace: ns, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) BuildClient() (kubernetes.Interface, error) {
	return cb.Fakeclient, nil

}

// OCClientBuilder is a mock openshift client bulder
type OCClientBuilder struct {
	token     string
	namespace string
	host      string
	Fake      *testing.Fake
}

func NewOCClientBuilder(token, namespace, host string, fake *testing.Fake) *OCClientBuilder {
	return &OCClientBuilder{
		token:     token,
		namespace: namespace,
		host:      host,
		Fake:      fake,
	}
}

func (cb *OCClientBuilder) WithToken(token string) mobile.OSClientBuilder {
	return &OCClientBuilder{token: token, host: cb.host, namespace: cb.namespace, Fake: cb.Fake}
}
func (cb *OCClientBuilder) WithNamespace(ns string) mobile.OSClientBuilder {
	return &OCClientBuilder{token: cb.token, host: cb.host, namespace: ns, Fake: cb.Fake}
}
func (cb *OCClientBuilder) WithHost(host string) mobile.OSClientBuilder {
	return &OCClientBuilder{token: cb.token, host: host, namespace: cb.namespace, Fake: cb.Fake}
}
func (cb *OCClientBuilder) WithHostAndNamespace(host, ns string) mobile.OSClientBuilder {
	return &OCClientBuilder{token: cb.token, host: host, namespace: ns, Fake: cb.Fake}
}
func (cb *OCClientBuilder) BuildClient() (client.Interface, error) {
	return testclient.NewClient(cb.host, cb.namespace, cb.token, cb.Fake), nil

}

type SCClientBuilder struct {
	Client mobile.SCCInterface
}

func (scc *SCClientBuilder) WithToken(token string) mobile.SCClientBuilder {
	return scc
}
func (scc *SCClientBuilder) WithHost(host string) mobile.SCClientBuilder {
	return scc
}
func (scc *SCClientBuilder) UseDefaultSAToken() mobile.SCClientBuilder {
	return scc
}
func (scc *SCClientBuilder) Build() (mobile.SCCInterface, error) {
	return scc.Client, nil
}

func NewSCClient() *SCClient {
	return &SCClient{Calls: map[string]int{}}
}

type SCClient struct {
	Err   error
	Calls map[string]int
}

func (sc *SCClient) BindToService(bindableService, targetSvcName string, bindingParams map[string]interface{}, bindableServiceNamespace, targetSvcNamespace string) error {
	sc.Calls["BindToService"]++
	return sc.Err
}
func (sc *SCClient) UnBindFromService(bindableService, targetSvcName, bindableServiceNamespace string) error {
	sc.Calls["UnBindFromService"]++
	return sc.Err
}
func (sc *SCClient) AddMobileApiKeys(targetSvcName, namespace string) error {
	sc.Calls["AddMobileApiKeys"]++
	return sc.Err
}
func (sc *SCClient) RemoveMobileApiKeys(targetSvcName, namespace string) error {
	sc.Calls["RemoveMobileApiKeys"]++
	return sc.Err
}
func (sc *SCClient) Called(f string) int {
	return sc.Calls[f]
}
