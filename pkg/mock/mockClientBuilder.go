package mock

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"k8s.io/client-go/kubernetes"
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
