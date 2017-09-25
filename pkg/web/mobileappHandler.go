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
	logger           *logrus.Logger
	appService       *app.Service
	appCruderBuilder mobile.AppRepoBuilder
}

// NewMobileAppHandler returns a new mobile app handler
func NewMobileAppHandler(logger *logrus.Logger, app mobile.AppRepoBuilder, appService *app.Service) *MobileAppHandler {
	return &MobileAppHandler{
		logger:           logger,
		appCruderBuilder: app,
		appService:       appService,
	}
}

// Read reads a mobileapp based on an id
func (m *MobileAppHandler) Read(rw http.ResponseWriter, req *http.Request) {
	// we return the actul request handleing function now that it has been configured.
	token := headers.DefaultTokenRetriever(req.Header)
	appRepo, err := m.appCruderBuilder.WithToken(token).Build()
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
	appRepo, err := m.appCruderBuilder.WithToken(token).Build()
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
	appRepo, err := m.appCruderBuilder.WithToken(token).Build()
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	params := mux.Vars(req)
	id := params["id"]

	if err := m.appService.Delete(appRepo, id); err != nil {
		err = errors.Wrap(err, "mobile app handler, failed to delete app")
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
}

// Create creates a mobileapp
func (m *MobileAppHandler) Create(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	appRepo, err := m.appCruderBuilder.WithToken(token).Build()
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

	rw.WriteHeader(http.StatusCreated)
}

// Update will update a mobile app
func (m *MobileAppHandler) Update(rw http.ResponseWriter, req *http.Request) {

}
