package clients

import (
	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

//Note a concern with this pattern is that it essentially becomes a registry that all handlers depend on, this will have an impact on tests
//as we will always first need to construct this before we can do anything

// TokenScopedClientBuilder builds a client bound to a particular token.
// if there is token passed it will attempt to use the default sa token
type TokenScopedClientBuilder struct {
	clientBuilder      mobile.ClientBuilder
	appRepoBuilder     mobile.AppRepoBuilder
	serviceRepoBuilder mobile.ServiceRepoBuilder
	namespace          string
	logger             *logrus.Logger
	mounterBuilder     mobile.MounterBuilder
	useSaToken         bool
	// this is initialised to the service acount token in the container
	SAToken string
}

// NewTokenScopedClientBuilder returns a new client builder that builds clients using the token provided
func NewTokenScopedClientBuilder(cb mobile.ClientBuilder, arb mobile.AppRepoBuilder, srv mobile.ServiceRepoBuilder, mb mobile.MounterBuilder, namespace string, logger *logrus.Logger) *TokenScopedClientBuilder {
	return &TokenScopedClientBuilder{
		clientBuilder:      cb,
		appRepoBuilder:     arb,
		serviceRepoBuilder: srv,
		namespace:          namespace,
		logger:             logger,
		mounterBuilder:     mb,
	}
}

func (rsb *TokenScopedClientBuilder) token(t string) string {
	if rsb.useSaToken {
		rsb.logger.Info("TokenScopedClientBuilder ignoring passed token and instead is using service account token for authentication")
		return rsb.SAToken
	}
	return t
}

//UseDefaultSAToken clones the client builder and sets it to use the service account token
func (rsb *TokenScopedClientBuilder) UseDefaultSAToken() mobile.TokenScopedClientBuilder {
	var cloned = *rsb
	cloned.useSaToken = true
	return &cloned
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

func (rsb *TokenScopedClientBuilder) VolumeMounterUnmounter(token string) (mobile.VolumeMounterUnmounter, error) {
	token = rsb.token(token)
	k8client, err := rsb.clientBuilder.WithToken(token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "client.rsb.VolumeMounterUnmounter -> failed to create request scoped kubernetes client with token")
	}
	return rsb.mounterBuilder.WithK8s(k8client).Build(), nil
}
