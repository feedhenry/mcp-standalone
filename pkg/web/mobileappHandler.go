package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/gorilla/mux"
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
	token := req.Header.Get(mobile.AuthHeader)
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
	token := req.Header.Get(mobile.AuthHeader)
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
	token := req.Header.Get(mobile.AuthHeader)
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
	token := req.Header.Get(mobile.AuthHeader)
	appRepo, err := m.tokenClientBuilder.MobileAppCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	uid := uuid.NewV4()
	decoder := json.NewDecoder(req.Body)
	app := &mobile.App{}
	app.APIKey = uid.String()
	if err := decoder.Decode(app); err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	if err := appRepo.Create(app); err != nil {
		handleCommonErrorCases(err, rw, m.logger)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// Update will update a mobile app
func (m *MobileAppHandler) Update(rw http.ResponseWriter, req *http.Request) {

}
