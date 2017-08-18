package mobile

import (
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
