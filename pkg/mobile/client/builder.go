package client

import (
	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

//TODO WE may want to move this out to the data package as it is sepcific to kubnernetes

// TokenScopedClientBuilder builds a client bound to a particular token.
// if there is token passed it will attempt to use the default sa token
type TokenScopedClientBuilder struct {
	clientBuilder      mobile.ClientBuilder
	appRepoBuilder     mobile.AppRepoBuilder
	serviceRepoBuilder mobile.ServiceRepoBuilder
	namespace          string
	logger             *logrus.Logger
	// this is initialised to the service acount token in the container
	SAToken string
}

//UseDefaultSAToken indicates to the builder to use the sa token
const UseDefaultSAToken = "satoken"

// NewTokenScopedClientBuilder returns a new client builder that builds clients using the token provided
func NewTokenScopedClientBuilder(cb mobile.ClientBuilder, arb mobile.AppRepoBuilder, srv mobile.ServiceRepoBuilder, namespace string, logger *logrus.Logger) *TokenScopedClientBuilder {
	return &TokenScopedClientBuilder{
		clientBuilder:      cb,
		appRepoBuilder:     arb,
		serviceRepoBuilder: srv,
		namespace:          namespace,
		logger:             logger,
	}
}

func (rsb *TokenScopedClientBuilder) token(t string) string {
	if t == UseDefaultSAToken {
		rsb.logger.Info("using service account token in client ")
		return rsb.SAToken
	}
	return t
}

// MobileAppCruder returns a token scoped MobileAppCruder
func (rsb *TokenScopedClientBuilder) MobileAppCruder(token string) (mobile.AppCruder, error) {
	token = rsb.token(token)
	k8s, err := rsb.K8s(token)
	if err != nil {
		return nil, err
	}
	return rsb.appRepoBuilder.WithClient(k8s.CoreV1().ConfigMaps(rsb.namespace)).Build(), nil
}

// K8s will build a token scoped kuberentes client
func (rsb *TokenScopedClientBuilder) K8s(token string) (kubernetes.Interface, error) {
	token = rsb.token(token)
	k8client, err := rsb.clientBuilder.WithToken(token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request scoped kubernetes client with token")

	}
	return k8client, nil
}

// MobileServiceCruder builds a token scoped service cruder
func (rsb *TokenScopedClientBuilder) MobileServiceCruder(token string) (mobile.ServiceCruder, error) {
	token = rsb.token(token)
	k8client, err := rsb.clientBuilder.WithToken(token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request scoped kubernetes client with token")
	}
	return rsb.serviceRepoBuilder.WithClient(k8client.CoreV1().Secrets(rsb.namespace)).Build(), nil

}
