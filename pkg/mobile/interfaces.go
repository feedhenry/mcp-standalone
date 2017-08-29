package mobile

import (
	"net/http"

	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type AppCruder interface {
	ReadByName(name string) (*App, error)
	Create(app *App) error
	DeleteByName(name string) error
	List() ([]*App, error)
	Update(app *App) (*App, error)
}

type ServiceCruder interface {
	List(AttrFilterFunc) ([]*Service, error)
	ListConfigs(AttrFilterFunc) ([]*ServiceConfig, error)
}

type Attributer interface {
	GetName() string
	GetLabels() map[string]string
}

type ClientBuilder interface {
	WithToken(token string) ClientBuilder
	WithNamespace(ns string) ClientBuilder
	WithHost(host string) ClientBuilder
	WithHostAndNamespace(host, ns string) ClientBuilder
	BuildClient() (kubernetes.Interface, error)
	BuildConfigMapClent() (corev1.ConfigMapInterface, error)
	BuildSecretClent() (corev1.SecretInterface, error)
}

type AppRepoBuilder interface {
	WithClient(c corev1.ConfigMapInterface) AppRepoBuilder
	Build() AppCruder
}

type ServiceRepoBuilder interface {
	WithClient(c corev1.SecretInterface) ServiceRepoBuilder
	Build() ServiceCruder
}

type TokenScopedClientBuilder interface {
	K8s(token string) (kubernetes.Interface, error)
	MobileAppCruder(token string) (AppCruder, error)
	MobileServiceCruder(token string) (ServiceCruder, error)
	UseDefaultSAToken() TokenScopedClientBuilder
}

type HTTPRequesterBuilder interface {
	Insecure(i bool) HTTPRequesterBuilder
	Timeout(t int) HTTPRequesterBuilder
	Build() ExternalHTTPRequester
}

type ExternalHTTPRequester interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}
