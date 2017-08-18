package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	kerror "k8s.io/apimachinery/pkg/api/errors"
)

// MobileAppHandler handle mobile actions
type MobileAppHandler struct {
	logger *logrus.Logger
}

// NewMobileAppHandler returns a new mobile app handler
func NewMobileAppHandler(logger *logrus.Logger) *MobileAppHandler {
	return &MobileAppHandler{
		logger: logger,
	}
}

// Read reads a mobileapp based on an id
func (m *MobileAppHandler) Read(appRepo mobile.AppCruder) http.HandlerFunc {
	// we return the actul request handleing function now that it has been configured.
	return func(rw http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		id := params["id"]
		encoder := json.NewEncoder(rw)
		if "" == id {
			http.Error(rw, "id cannot be empty", http.StatusBadRequest)
			return
		}
		app, err := appRepo.ReadByName(id)
		if err != nil {
			m.handlerError(err, rw)
			return
		}
		if err := encoder.Encode(app); err != nil {
			m.handlerError(err, rw)
			return
		}
	}
}

// List will list mobile apps
func (m *MobileAppHandler) List(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		apps, err := appRepo.List()
		if err != nil {
			m.handlerError(err, rw)
			return
		}
		encoder := json.NewEncoder(rw)
		if err := encoder.Encode(apps); err != nil {
			m.handlerError(err, rw)
			return
		}
	}
}

// Delete will delete a mobile app
func (m *MobileAppHandler) Delete(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		id := params["id"]
		if err := appRepo.DeleteByName(id); err != nil {
			m.handlerError(err, rw)
			return
		}
	}
}

// Create creates a mobileapp
func (m *MobileAppHandler) Create(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		app := &mobile.App{}
		if err := decoder.Decode(app); err != nil {
			m.handlerError(err, rw)
			return
		}
		//add validation
		if err := appRepo.Create(app); err != nil {
			m.handlerError(err, rw)
			return
		}
		rw.WriteHeader(http.StatusCreated)

	}
}

// Update will update a mobile app
func (m *MobileAppHandler) Update(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
	}
}

// likely we will eventually abstract out
func (m *MobileAppHandler) handlerError(err error, rw http.ResponseWriter) {
	err = errors.Cause(err)
	if mobile.IsNotFoundError(err) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	if e, ok := err.(*mobile.StatusError); ok {
		m.logger.Error(fmt.Sprintf("status error occurred %+v", err))
		http.Error(rw, err.Error(), e.StatusCode())
		return
	}
	if e, ok := err.(*kerror.StatusError); ok {
		m.logger.Error(fmt.Sprintf("status error occurred %+v", err))
		http.Error(rw, e.Error(), int(e.Status().Code))
		return
	}
	m.logger.Error(fmt.Sprintf("unexpected error occurred %+v", err))
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}
