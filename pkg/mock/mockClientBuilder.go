package mock

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ClientBuilder struct {
	token      string
	namespace  string
	host       string
	Fakeclient kubernetes.Interface
}

func (cb *ClientBuilder) WithToken(token string) mobile.ClientBuilder {
	return &ClientBuilder{token: token, host: cb.host, namespace: cb.namespace, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) WithNamespace(ns string) mobile.ClientBuilder {
	return &ClientBuilder{token: cb.token, host: cb.host, namespace: ns, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) WithHost(host string) mobile.ClientBuilder {
	return &ClientBuilder{token: cb.token, host: host, namespace: cb.namespace, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) WithHostAndNamespace(host, ns string) mobile.ClientBuilder {
	return &ClientBuilder{token: cb.token, host: host, namespace: ns, Fakeclient: cb.Fakeclient}
}
func (cb *ClientBuilder) BuildClient() (kubernetes.Interface, error) {
	return cb.Fakeclient, nil

}
func (cb *ClientBuilder) BuildConfigMapClent() (corev1.ConfigMapInterface, error) {
	return cb.Fakeclient.CoreV1().ConfigMaps(cb.namespace), nil
}
func (cb *ClientBuilder) BuildSecretClent() (corev1.SecretInterface, error) {
	return cb.Fakeclient.CoreV1().Secrets(cb.namespace), nil
}

type AppRepoBuilder struct {
	AppCruder *MockAppCruder
}

func (arb *AppRepoBuilder) WithClient(c corev1.ConfigMapInterface) mobile.AppRepoBuilder {
	return &AppRepoBuilder{AppCruder: arb.AppCruder}
}
func (arb *AppRepoBuilder) Build() mobile.AppCruder {
	return arb.AppCruder
}

type MockAppCruder struct {
	App  *mobile.App
	Apps []*mobile.App
	Err  error
}

func (mac *MockAppCruder) ReadByName(name string) (*mobile.App, error) {

	return mac.App, mac.Err
}
func (mac *MockAppCruder) Create(app *mobile.App) error {
	return mac.Err
}
func (mac *MockAppCruder) DeleteByName(name string) error {
	return mac.Err
}
func (mac *MockAppCruder) List() ([]*mobile.App, error) {
	return mac.Apps, mac.Err
}
func (mac *MockAppCruder) Update(app *mobile.App) (*mobile.App, error) {
	return mac.App, mac.Err
}
