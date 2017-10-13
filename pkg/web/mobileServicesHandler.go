package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/openshift"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// MobileServiceHandler handles endpoints associated with mobile enabled services. It will list services in the namespace that
// it knows about so that they can be rendered in the MCP
type MobileServiceHandler struct {
	logger                   *logrus.Logger
	mobileIntegrationService *integration.MobileService
	mounterBuilder           mobile.MounterBuilder
	serviceRepoBuilder       mobile.ServiceRepoBuilder
	metricsGetter            mobile.MetricsGetter
	userRepoBuilder          mobile.UserRepoBuilder
	authCheckerBuilder       mobile.AuthCheckerBuilder
	sccClientBuilder         mobile.SCClientBuilder
}

// NewMobileServiceHandler returns a new MobileServiceHandler
func NewMobileServiceHandler(logger *logrus.Logger, integrationService *integration.MobileService, mounterBuilder mobile.MounterBuilder,
	mg mobile.MetricsGetter, serviceRepoBuilder mobile.ServiceRepoBuilder, userRepoBuilder mobile.UserRepoBuilder, authCheckerBuilder mobile.AuthCheckerBuilder,
	sccClientBuilder mobile.SCClientBuilder) *MobileServiceHandler {
	return &MobileServiceHandler{
		logger: logger,
		mobileIntegrationService: integrationService,
		metricsGetter:            mg,
		mounterBuilder:           mounterBuilder,
		serviceRepoBuilder:       serviceRepoBuilder,
		userRepoBuilder:          userRepoBuilder,
		authCheckerBuilder:       authCheckerBuilder,
		sccClientBuilder:         sccClientBuilder,
	}
}

// List allows you to list mobile services
func (msh *MobileServiceHandler) List(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	serviceCruder, err := msh.serviceRepoBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	userRepo := msh.userRepoBuilder.WithToken(token).WithClient(&openshift.UserAccess{}).Build()
	authChecker := msh.authCheckerBuilder.WithToken(token).WithUserRepo(userRepo).IgnoreCerts().Build()
	client := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()
	svc, err := msh.mobileIntegrationService.DiscoverMobileServices(serviceCruder, authChecker, client)
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
	serviceCruder, err := msh.serviceRepoBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	userRepo := msh.userRepoBuilder.WithToken(token).WithClient(&openshift.UserAccess{}).Build()
	authChecker := msh.authCheckerBuilder.WithToken(token).WithUserRepo(userRepo).IgnoreCerts().Build()
	client := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()

	if withIntegrations != "" {
		ms, err = msh.mobileIntegrationService.ReadMobileServiceAndIntegrations(serviceCruder, authChecker, serviceName, client)
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
	serviceCruder, err := msh.serviceRepoBuilder.WithToken(token).Build()
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
	clientServiceName := strings.ToLower(params["clientService"])
	serviceName := strings.ToLower(params["targetService"])
	if len(serviceName) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure: provided targetService must not be empty"), rw, msh.logger)
		return
	}
	if len(clientServiceName) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure: provided clientServiceName must not be empty"), rw, msh.logger)
		return
	}

	scClient, err := msh.sccClientBuilder.WithToken(token).Build()
	if err != nil {
		err = errors.Wrap(err, "web.msh.Configure: failed to create the service catalog client")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}

	serviceCruder, err := msh.serviceRepoBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Configure: could not create service cruder"), rw, msh.logger)
		return
	}

	if err := msh.mobileIntegrationService.BindService(scClient, serviceCruder, clientServiceName, serviceName); err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Configure: could not create binding for service : '"+serviceName+"' for target: '"+clientServiceName), rw, msh.logger)
		return
	}
}

// Deconfigure removes configuration for components binding
func (msh *MobileServiceHandler) Deconfigure(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := headers.DefaultTokenRetriever(req.Header)

	params := mux.Vars(req)
	componentType := strings.ToLower(params["componentType"])
	componentName := strings.ToLower(params["componentName"])
	secret := params["secret"]
	if len(componentType) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided componentType must not be empty"), rw, msh.logger)
		return
	}
	if len(componentName) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided componentName must not be empty"), rw, msh.logger)
		return
	}
	if len(secret) == 0 {
		handleCommonErrorCases(errors.New("web.msh.Configure -> provided secret must not be empty"), rw, msh.logger)
		return
	}

	serviceCruder, err := msh.serviceRepoBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not create service cruder"), rw, msh.logger)
		return
	}

	unmounter, err := msh.mounterBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not create volume unmounter"), rw, msh.logger)
		return
	}

	err = msh.mobileIntegrationService.UnmountSecretInComponent(serviceCruder, unmounter, componentType, componentName, secret)
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not unmount secret: '"+secret+"' from component: '"+componentType+":"+componentName+"'"), rw, msh.logger)
		return
	}
	return
}

func (msh *MobileServiceHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	token := headers.DefaultTokenRetriever(req.Header)
	params := mux.Vars(req)
	serviceCruder, err := msh.serviceRepoBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(errors.Wrap(err, "web.msh.Deconfigure -> could not create service cruder"), rw, msh.logger)
		return
	}
	serviceName := params["name"]
	if serviceName == "" {
		http.Error(rw, "service name cannot be empty", http.StatusBadRequest)
		return
	}
	if err := serviceCruder.Delete(serviceName); err != nil {
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

	if strings.HasPrefix(serviceName, "fh-sync-server") {
		serviceName = "fh-sync-server"
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
