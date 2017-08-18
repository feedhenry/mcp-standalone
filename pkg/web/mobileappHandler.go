package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/gorilla/mux"
)

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

func (m *MobileAppHandler) Delete(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
	}
}

func (m *MobileAppHandler) Create(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		m.logger.Debug("create called")
		decoder := json.NewDecoder(req.Body)
		app := &mobile.App{}
		if err := decoder.Decode(app); err != nil {
			m.handlerError(err, rw)
			return
		}
		m.logger.Debug(" decoded app ", app)
		//add validation
		if err := appRepo.Create(app); err != nil {
			m.handlerError(err, rw)
			return
		}
		m.logger.Debug(" created ", app)
	}
}

func (m *MobileAppHandler) Update(appRepo mobile.AppCruder) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
	}
}

// likely we will eventually abstract out
func (m *MobileAppHandler) handlerError(err error, rw http.ResponseWriter) {
	if mobile.IsNotFoundError(err) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	m.logger.Error(fmt.Sprintf("unexpected error occurred %+v", err))
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}
