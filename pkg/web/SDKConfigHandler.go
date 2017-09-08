package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// SDKConfigHandler handles sdk configuration requests
type SDKConfigHandler struct {
	mobileIntegrationService *integration.MobileService
	tokenScopedBuilder       mobile.TokenScopedClientBuilder
	logger                   *logrus.Logger
}

// NewSDKConfigHandler returns an sdk handler
func NewSDKConfigHandler(logger *logrus.Logger, service *integration.MobileService, builder mobile.TokenScopedClientBuilder) *SDKConfigHandler {
	return &SDKConfigHandler{
		mobileIntegrationService: service,
		logger:             logger,
		tokenScopedBuilder: builder,
	}
}

func (sdk *SDKConfigHandler) Read(rw http.ResponseWriter, req *http.Request) {
	//need to read the mobile app and authenticate its apikey. There wont be an openshift token
	apiKey := req.Header.Get(mobile.AppAPIKeyHeader)
	params := mux.Vars(req)
	id := params["id"]
	if apiKey == "" {
		http.Error(rw, "missing api key", 401)
	}
	//TODO maybe bring this  apiKey check out of this handler
	//need to use the serviceaccount token here to read and check the app key and svcs
	appCruder, err := sdk.tokenScopedBuilder.UseDefaultSAToken().MobileAppCruder("")
	if err != nil {
		err = errors.Wrap(err, "failed to setup mobile app cruder using sa token")
		handleCommonErrorCases(err, rw, sdk.logger)
		return
	}
	svcCruder, err := sdk.tokenScopedBuilder.UseDefaultSAToken().MobileServiceCruder("")
	if err != nil {
		err = errors.Wrap(err, "failed to create token scoped service client")
		handleCommonErrorCases(err, rw, sdk.logger)
		return
	}
	//before returning any information check the passed api key is the same as the app objects generated key.
	app, err := appCruder.ReadByName(id)
	if err != nil {
		handleCommonErrorCases(err, rw, sdk.logger)
		return
	}
	if apiKey != app.APIKey {
		http.Error(rw, "unauthorised ", http.StatusUnauthorized)
		return
	}
	configs, err := sdk.mobileIntegrationService.GenerateMobileServiceConfigs(svcCruder)
	if err != nil {
		handleCommonErrorCases(err, rw, sdk.logger)
		return
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(configs); err != nil {
		handleCommonErrorCases(err, rw, sdk.logger)
		return
	}
}
