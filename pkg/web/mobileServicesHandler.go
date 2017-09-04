package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/pkg/errors"
)

// MobileServiceHandler handles endpoints associated with mobile enabled services. It will list services in the namespace that
// it knows about so that they can be rendered in the MCP
type MobileServiceHandler struct {
	logger                   *logrus.Logger
	mobileIntegrationService *integration.MobileService
	tokenClientBuilder       mobile.TokenScopedClientBuilder
}

// NewMobileServiceHandler returns a new MobileServiceHandler
func NewMobileServiceHandler(logger *logrus.Logger, integrationService *integration.MobileService, tokenClientBuilder mobile.TokenScopedClientBuilder) *MobileServiceHandler {
	return &MobileServiceHandler{
		logger: logger,
		mobileIntegrationService: integrationService,
		tokenClientBuilder:       tokenClientBuilder,
	}
}

type mobileServiceConfigurationRequest struct {
	Component string `json:"component"`
	Service   string `json:"service"`
	Namespace string `json:"namespace"`
}

// List allows you to list mobile services
func (msh *MobileServiceHandler) List(rw http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(mobile.AuthHeader)
	serviceCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	svc, err := msh.mobileIntegrationService.DiscoverMobileServices(serviceCruder)
	if err != nil {
		err = errors.Wrap(err, "attempted to list mobile services")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(svc); err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
}

// Configure configures components binding
func (msh *MobileServiceHandler) Configure(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := req.Header.Get(mobile.AuthHeader)

	decoder := json.NewDecoder(req.Body)
	var conf mobileServiceConfigurationRequest
	err := decoder.Decode(&conf)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	conf.Component = strings.ToLower(conf.Component)
	conf.Service = strings.ToLower(conf.Service)
	// TODO move this out of the handler
	serviceCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	services, err := msh.mobileIntegrationService.FindByNames([]string{conf.Service}, serviceCruder)
	if err != nil {
		err = errors.Wrap(err, "attempted to list mobile services")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	if len(services) == 0 {
		err = errors.New("service to configure not found")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	components, err := msh.mobileIntegrationService.FindByNames([]string{conf.Component}, serviceCruder)
	if err != nil {
		err = errors.Wrap(err, "attempted to list mobile services")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	if len(components) == 0 {
		err = errors.New("component to configure not found")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	k8sClient, err := msh.tokenClientBuilder.K8s(token)

	deploy, err := msh.mobileIntegrationService.MountSecretForComponent(k8sClient, services[0].BindingSecretName, conf.Component, conf.Service, conf.Namespace)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(deploy); err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
}
