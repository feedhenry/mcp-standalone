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
	metricsGetter            mobile.MetricsGetter
}

// NewMobileServiceHandler returns a new MobileServiceHandler
func NewMobileServiceHandler(logger *logrus.Logger, integrationService *integration.MobileService, tokenClientBuilder mobile.TokenScopedClientBuilder, mg mobile.MetricsGetter) *MobileServiceHandler {
	return &MobileServiceHandler{
		logger: logger,
		mobileIntegrationService: integrationService,
		tokenClientBuilder:       tokenClientBuilder,
		metricsGetter:            mg,
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

// Read details of a specific service provided in the name URL parameter
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
		ms, err = msh.mobileIntegrationService.ReadMobileServiceAndIntegrations(serviceCruder, serviceName)
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

// Configure configures components binding
func (msh *MobileServiceHandler) Configure(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := headers.DefaultTokenRetriever(req.Header)

	params := mux.Vars(req)
	component := strings.ToLower(params["component"])
	secret := strings.ToLower(params["secret"])
	if len(component) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided component must not be empty"), rw, msh.logger)
		return
	}
	if len(secret) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided secret must not be empty"), rw, msh.logger)
		return
	}

	mounter, err := msh.tokenClientBuilder.VolumeMounterUnmounter(token)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Configure -> could not create mounter"), rw, msh.logger)
		return
	}

	svcCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Configure -> could not create service cruder"), rw, msh.logger)
		return
	}

	err = msh.mobileIntegrationService.MountSecretForComponent(svcCruder, mounter, component, secret)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Configure -> could not mount secret: '"+secret+"' into component: '"+component+"'"), rw, msh.logger)
		return
	}

	return
}

// Deconfigure removes configuration for components binding
func (msh *MobileServiceHandler) Deconfigure(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := headers.DefaultTokenRetriever(req.Header)

	params := mux.Vars(req)
	component := params["component"]
	secret := params["secret"]
	if len(component) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided component must not be empty"), rw, msh.logger)
		return
	}
	if len(secret) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided secret must not be empty"), rw, msh.logger)
		return
	}

	svcCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not create service cruder"), rw, msh.logger)
		return
	}

	unmounter, err := msh.tokenClientBuilder.VolumeMounterUnmounter(token)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not create volume unmounter"), rw, msh.logger)
		return
	}

	err = msh.mobileIntegrationService.UnmountSecretInComponent(svcCruder, unmounter, component, secret)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not unmount secret: '"+secret+"' from component: '"+component+"'"), rw, msh.logger)
		return
	}
	return
}

func (msh *MobileServiceHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := headers.DefaultTokenRetriever(req.Header)
	params := mux.Vars(req)
	svcCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not create service cruder"), rw, msh.logger)
		return
	}
	serviceName := params["name"]
	if serviceName == "" {
		http.Error(rw, "service name cannot be empty", http.StatusBadRequest)
		return
	}
	if err := svcCruder.Delete(serviceName); err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Delete could not delete service "), rw, msh.logger)
		return
	}

}

func (msh *MobileServiceHandler) GetMetrics(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	serviceName := params["name"]
	metric := req.URL.Query().Get("metric")
	var encoder = json.NewEncoder(rw)
	if serviceName == "" {
		http.Error(rw, "missing service name. cannot be empty", http.StatusBadRequest)
		return
	}

	if metric == "" {
		metrics := msh.metricsGetter.GetAll(serviceName)
		if nil == metrics {
			http.Error(rw, "no metrics found for service", http.StatusNotFound)
			return
		}
		if err := encoder.Encode(metrics); err != nil {
			handleCommonErrorCases(err, rw, msh.logger)
			return
		}
		return
	}
	serviceMetric := msh.metricsGetter.GetOne(serviceName, metric)
	if nil == serviceMetric {
		http.Error(rw, "no metrics or metric found for service", http.StatusNotFound)
		return
	}
	if err := encoder.Encode(serviceMetric); err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
}
