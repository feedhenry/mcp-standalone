package middleware

import (
	"net/http"

	"github.com/feedhenry/mobile-server/pkg/mobile"
	"k8s.io/client-go/kubernetes"
)

// Builder is a piece of middleware that is used to build clients based on the details in the request.
type Builder struct {
	clientBuilder  mobile.ClientBuilder
	appRepoBuilder mobile.AppRepoBuilder
	namespace      string
}

func NewBuilder(cb mobile.ClientBuilder, arb mobile.AppRepoBuilder, namespace string) *Builder {
	return &Builder{
		clientBuilder:  cb,
		appRepoBuilder: arb,
		namespace:      namespace,
	}
}

// K8ClientHTTPHandlerFunc is a func that wants a configured kubernetes client to use and returns a http request handling func to complete the request
type K8ClientHTTPHandlerFunc func(k8client kubernetes.Interface) http.HandlerFunc

// HandleKubernetesClient accepts the incoming request and a K8ClientHTTPHandlerFunc instanciates a kubernetes client and configures the K8ClientHTTPHandlerFunc
//with that client before handing the request off
func (c *Builder) HandleKubernetesClient(handler K8ClientHTTPHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("auth")
		k8client, err := c.clientBuilder.WithToken(token).BuildClient()
		if err != nil {
			return
		}
		configuredHandler := handler(k8client)
		configuredHandler(rw, req)
	}
}

// RepoHTTPHandlerFunc is a func that wants a configured kubernetes client to use and returns a http request handling func to complete the request
type RepoHTTPHandlerFunc func(dataRepo mobile.AppCruder) http.HandlerFunc

// HandleRepo  accepts an incoming request and a function that want a app repo
func (c *Builder) HandleRepo(handler RepoHTTPHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("x-auth")
		k8client, err := c.clientBuilder.WithToken(token).BuildClient()
		if err != nil {
			return
		}
		repo := c.appRepoBuilder.WithClient(k8client.CoreV1().ConfigMaps(c.namespace)).Build()
		configuredHandler := handler(repo)
		configuredHandler(rw, req)
	}
}
