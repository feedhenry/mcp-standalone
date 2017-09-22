package mobile

import (
	"net/http"

	"k8s.io/client-go/kubernetes"
)

type AppCruder interface {
	ReadByName(name string) (*App, error)
	Create(app *App) error
	DeleteByName(name string) error
	List() ([]*App, error)
	Update(app *App) (*App, error)
	UpdateAppAPIKeys(app *App) error
	CreateAppAPIKeyMap() error
	RemoveAppAPIKeyByID(appID string) error
}

type ServiceCruder interface {
	List(AttrFilterFunc) ([]*Service, error)
	Read(name string) (*Service, error)
	ListConfigs(AttrFilterFunc) ([]*ServiceConfig, error)
	UpdateEnabledIntegrations(svcName string, integrations map[string]string) error
	Create(ms *Service) error
	Delete(serviceID string) error
}

type Attributer interface {
	GetName() string
	GetLabels() map[string]string
	GetType() string
}

//TODO probably not a core interface but rather we should wrap it inside the other repos as a dependency and have it consumed via the builders
type ClientBuilder interface {
	WithToken(token string) ClientBuilder
	WithNamespace(ns string) ClientBuilder
	WithHost(host string) ClientBuilder
	WithHostAndNamespace(host, ns string) ClientBuilder
	BuildClient() (kubernetes.Interface, error)
}

type AppRepoBuilder interface {
	WithToken(token string) AppRepoBuilder
	//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
	UseDefaultSAToken() AppRepoBuilder
	Build() (AppCruder, error)
}

type UserRepoBuilder interface {
	WithToken(token string) UserRepoBuilder
	WithClient(client UserAccessChecker) UserRepoBuilder
	Build() UserRepo
}

type UserRepo interface {
	GetUser() (*User, error)
}

// TODO prob can remote the WithClient and instead use NewRepoBuilder(c corev1.ConfigMapInterface) and have this just expose Build() and perhaps add WithToken(token string)
type ServiceRepoBuilder interface {
	WithToken(token string) ServiceRepoBuilder
	//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
	UseDefaultSAToken() ServiceRepoBuilder
	Build() (ServiceCruder, error)
}

type TokenScopedClientBuilder interface {
	K8s(token string) (kubernetes.Interface, error)
	UseDefaultSAToken() TokenScopedClientBuilder
	VolumeMounterUnmounter(token string) (VolumeMounterUnmounter, error)
	AuthChecker(token string, ignoreCerts bool) AuthChecker
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

// MounterBuilder creates VolumeMounterUnmounter objects
type MounterBuilder interface {
	Build() (VolumeMounterUnmounter, error)
	WithToken(token string) MounterBuilder
	//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
	UseDefaultSAToken() MounterBuilder
}

// VolumeMounter defines an interface for mounting volumes into services
type VolumeMounter interface {
	Mount(service, clientService *Service) error
}

// VolumeUnmounter defines an interface for unmounting volumes mounted in services
type VolumeUnmounter interface {
	Unmount(service, clientService *Service) error
}

// VolumeMounterUnmounter can both mount and unmount volumes
type VolumeMounterUnmounter interface {
	VolumeMounter
	VolumeUnmounter
}

// AuthCheckerBuilder builds AuthCheckers
type AuthCheckerBuilder interface {
	Build() AuthChecker
	WithToken(token string) AuthCheckerBuilder
	WithUserRepo(repo UserRepo) AuthCheckerBuilder
	IgnoreCerts() AuthCheckerBuilder
}

// AuthChecker performs a check for authorization to write the provided resource in the provided namespace
type AuthChecker interface {
	Check(resource, namespace string, client ExternalHTTPRequester) (bool, error)
}

type UserAccessChecker interface {
	ReadUserFromToken(host, token string, insecure bool) (*User, error)
}

type MetricsGetter interface {
	GetAll(serviceName string) []*GatheredMetric
	GetOne(serviceName, metric string) *GatheredMetric
}
