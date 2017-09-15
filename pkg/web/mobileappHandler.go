package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// MobileAppHandler handle mobile actions
type MobileAppHandler struct {
	logger             *logrus.Logger
	tokenClientBuilder mobile.TokenScopedClientBuilder
	appService         *app.Service
}

// NewMobileAppHandler returns a new mobile app handler
func NewMobileAppHandler(logger *logrus.Logger, tokenClientBuilder mobile.TokenScopedClientBuilder, appService *app.Service) *MobileAppHandler {
	return &MobileAppHandler{
		logger:             logger,
		tokenClientBuilder: tokenClientBuilder,
		appService:         appService,
	}
}

// Read reads a mobileapp based on an id
func (m *MobileAppHandler) Read(rw http.ResponseWriter, req *http.Request) {
	// we return the actul request handleing function now that it has been configured.
	token := headers.DefaultTokenRetriever(req.Header)
	appRepo, err := m.tokenClientBuilder.MobileAppCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	params := mux.Vars(req)
	id := params["id"]
	encoder := json.NewEncoder(rw)
	if "" == id {
		http.Error(rw, "id cannot be empty", http.StatusBadRequest)
		return
	}
	app, err := appRepo.ReadByName(id)
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	if err := encoder.Encode(app); err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
}

// List will list mobile apps
func (m *MobileAppHandler) List(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	appRepo, err := m.tokenClientBuilder.MobileAppCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	apps, err := appRepo.List()
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(apps); err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}

}

// Delete will delete a mobile app
func (m *MobileAppHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	appRepo, err := m.tokenClientBuilder.MobileAppCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	params := mux.Vars(req)
	id := params["id"]
	if err := appRepo.DeleteByName(id); err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}

}

// Create creates a mobileapp
func (m *MobileAppHandler) Create(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	appRepo, err := m.tokenClientBuilder.MobileAppCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	app := &mobile.App{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(app); err != nil {
		err = errors.Wrap(err, "mobile app create: Attempted to decode payload in app")
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	app.MetaData = map[string]string{}

	if err := m.appService.Create(appRepo, app); err != nil {
		err = errors.Wrap(err, "mobile app handler, failed to create app")
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	//todo logic is creeping in here should only be for parsing and rendering. Move to mobile package

	rw.WriteHeader(http.StatusCreated)
}

// Update will update a mobile app
func (m *MobileAppHandler) Update(rw http.ResponseWriter, req *http.Request) {

}
