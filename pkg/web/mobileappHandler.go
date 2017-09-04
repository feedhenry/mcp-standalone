package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// MobileAppHandler handle mobile actions
type MobileAppHandler struct {
	logger             *logrus.Logger
	tokenClientBuilder mobile.TokenScopedClientBuilder
}

// NewMobileAppHandler returns a new mobile app handler
func NewMobileAppHandler(logger *logrus.Logger, tokenClientBuilder mobile.TokenScopedClientBuilder) *MobileAppHandler {
	return &MobileAppHandler{
		logger:             logger,
		tokenClientBuilder: tokenClientBuilder,
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
	uid := uuid.NewV4()
	decoder := json.NewDecoder(req.Body)

	app := &mobile.App{MetaData: map[string]string{}}

	if err := decoder.Decode(app); err != nil {
		err = errors.Wrap(err, "mobile app create: Attempted to decode payload in app")
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	//todo logic is creeping in here should only be for parsing and rendering. Move to mobile package
	app.APIKey = uid.String()
	switch app.ClientType {
	case "android":
		app.MetaData["icon"] = "fa-android"
		break
	case "iOS":
		app.MetaData["icon"] = "fa-apple"
		break
	case "cordova":
		app.MetaData["icon"] = "icon-cordova"
		break
	}

	if err := appRepo.Create(app); err != nil {
		err = errors.Wrap(err, "mobile app create: Attempted to create app via app repo")
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// Update will update a mobile app
func (m *MobileAppHandler) Update(rw http.ResponseWriter, req *http.Request) {

}
