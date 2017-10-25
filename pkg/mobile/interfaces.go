package mobile

import (
	"net/http"

	"io"

	"net/url"

	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"k8s.io/client-go/kubernetes"
)

type AppCruder interface {
	ReadByName(name string) (*App, error)
	Create(app *App) error
	DeleteByName(name string) error
	List() ([]*App, error)
	Update(app *App) (*App, error)
	CreateAPIKeyMap() error
	AddAPIKeyToMap(app *App) error
	RemoveAPIKeyFromMap(appID string) error
}

type ServiceCruder interface {
	List(AttrFilterFunc) ([]*Service, error)
	Read(id string) (*Service, error)
	ListConfigs(AttrFilterFunc) ([]*ServiceConfig, error)
	UpdateEnabledIntegrations(svcName string, integrations map[string]string) error
	Create(ms *Service) error
	Delete(serviceID string) error
}

type BuildCruder interface {
	Create(b *BuildConfig) error
	AddBuildAsset(asset BuildAsset) (string, error)
	AddDownload(buildName string, dl *BuildDownload) error
	GetDownload(buildName string) (*BuildDownload, error)
	Status(buildName string) (*BuildStatus, error)
}

type ArtifactRetriever interface {
	Retrieve(location *url.URL, token string) (io.ReadCloser, error)
}

type Attributer interface {
	GetName() string
	GetLabels() map[string]string
	GetType() string
}

//TODO probably not a core interface but rather we should wrap it inside the other repos as a dependency and have it consumed via the builders
type K8ClientBuilder interface {
	WithToken(token string) K8ClientBuilder
	WithNamespace(ns string) K8ClientBuilder
	WithHost(host string) K8ClientBuilder
	WithHostAndNamespace(host, ns string) K8ClientBuilder
	BuildClient() (kubernetes.Interface, error)
}

type SCClientBuilder interface {
	WithToken(token string) SCClientBuilder
	WithHost(host string) SCClientBuilder
	UseDefaultSAToken() SCClientBuilder
	Build() (SCCInterface, error)
}

type SCCInterface interface {
	BindToService(bindableService, targetSvcName string, bindingParams map[string]interface{}, bindableServiceNamespace, targetSvcNamespace string) error
	UnBindFromService(bindableService, targetSvcName, bindableServiceNamespace string) error
	AddMobileApiKeys(targetSvcName, namespace string) error
	RemoveMobileApiKeys(targetSvcName, namespace string) error
}

type OSClientBuilder interface {
	WithToken(token string) OSClientBuilder
	WithNamespace(ns string) OSClientBuilder
	WithHost(host string) OSClientBuilder
	WithHostAndNamespace(host, ns string) OSClientBuilder
	//TODO change to Build to be consistent
	BuildClient() (client.Interface, error)
}

type AppRepoBuilder interface {
	WithToken(token string) AppRepoBuilder
	//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
	UseDefaultSAToken() AppRepoBuilder
	Build() (AppCruder, error)
}

type BuildRepoBuilder interface {
	WithToken(token string) BuildRepoBuilder
	//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
	UseDefaultSAToken() BuildRepoBuilder
	Build() (BuildCruder, error)
}

type UserRepoBuilder interface {
	WithToken(token string) UserRepoBuilder
	WithClient(client UserAccessChecker) UserRepoBuilder
	Build() UserRepo
}

type UserRepo interface {
	GetUser() (*User, error)
}

type ServiceRepoBuilder interface {
	WithToken(token string) ServiceRepoBuilder
	//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
	UseDefaultSAToken() ServiceRepoBuilder
	Build() (ServiceCruder, error)
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
