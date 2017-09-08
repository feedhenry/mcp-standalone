package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/gorilla/mux"
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

// List allows you to list mobile services
func (msh *MobileServiceHandler) List(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
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

func (msh *MobileServiceHandler) Read(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	params := mux.Vars(req)
	serviceName := params["name"]
	withIntegrations := req.URL.Query().Get("withIntegrations")
	var ms *mobile.Service
	var err error
	if serviceName == "" {
		http.Error(rw, "service name cannot be empty ", http.StatusBadRequest)
		return
	}
	serviceCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	if withIntegrations != "" {
		fmt.Println("with Integrations", serviceName)
		ms, err = msh.mobileIntegrationService.ReadMoileServiceAndIntegrations(serviceCruder, serviceName)
		if err != nil {
			handleCommonErrorCases(err, rw, msh.logger)
			return
		}
	} else {
		ms, err = serviceCruder.Read(serviceName)
		if err != nil {
			err = errors.Wrap(err, "MobileServiceHandler : failed to read mobile service ")
			handleCommonErrorCases(err, rw, msh.logger)
			return
		}
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(ms); err != nil {
		err = errors.Wrap(err, "MobileServiceHandler: failed to encode response")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
}

// Create will create a mobile service
func (msh *MobileServiceHandler) Create(rw http.ResponseWriter, req *http.Request) {
	ms := mobile.NewMobileService()
	token := headers.DefaultTokenRetriever(req.Header)
	serviceCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		err = errors.Wrap(err, "failed to setup service cruder based on token")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(ms); err != nil {
		err = errors.Wrap(err, "failed to decode mobile app ")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	if err := serviceCruder.Create(ms); err != nil {
		err = errors.Wrap(err, "failed to create mobile app")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// Configure configures components binding TODO NEEDS A REFACTOR
func (msh *MobileServiceHandler) Configure(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := headers.DefaultTokenRetriever(req.Header)

	decoder := json.NewDecoder(req.Body)
	var conf mobile.ServiceIntegration
	err := decoder.Decode(&conf)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	conf.Component = strings.ToLower(conf.Component)
	conf.Service = strings.ToLower(conf.Service)

	// TODO move this out of the handler

	k8sClient, err := msh.tokenClientBuilder.K8s(token)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	svcCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	deploy, err := msh.mobileIntegrationService.MountSecretForComponent(svcCruder, k8sClient, conf.Service, conf.Component, conf.Service, conf.Namespace, conf.ComponentSecret)
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
