package web

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/feedhenry/mobile-server/pkg/web/middleware"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	kerror "k8s.io/apimachinery/pkg/api/errors"
)

// NewRouter sets up the HTTP Router
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	return r
}

// BuildHTTPHandler puts together our HTTPHandler
func BuildHTTPHandler(r *mux.Router, access *middleware.Access) http.Handler {
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n := negroni.New(recovery)
	cors := middleware.Cors{}
	n.UseFunc(cors.Handle)
	if access != nil {

		n.UseFunc(access.Handle)
	} else {
		fmt.Println("access control is turned off ")
	}
	n.UseHandler(r)
	return n
}

// MobileAppRoute configure and setup the /mobileapp route. The middleware.Builder is responsible for building per request instances of clients
func MobileAppRoute(r *mux.Router, handler *MobileAppHandler) {
	r.HandleFunc("/mobileapp", handler.Create).Methods("POST")
	r.HandleFunc("/mobileapp/{id}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/mobileapp/{id}", handler.Read).Methods("GET")
	r.HandleFunc("/mobileapp", handler.List).Methods("GET")
	r.HandleFunc("/mobileapp/{id}", handler.Update).Methods("PUT")
}

//SDKConfigRoute configures and sets up the /sdk routes
func SDKConfigRoute(r *mux.Router, handler *SDKConfigHandler) {
	r.HandleFunc("/sdk/mobileapp/{id}/config", handler.Read).Methods("GET")
}

// SysRoute congifures and sets up the /sys/* route
func SysRoute(r *mux.Router, handler *SysHandler) {
	r.HandleFunc("/sys/info/ping", handler.Ping).Methods("GET")
	r.HandleFunc("/sys/info/health", handler.Health).Methods("GET")
}

// MobileServiceRoute configures and sets up the /mobileservice routes
func MobileServiceRoute(r *mux.Router, handler *MobileServiceHandler) {
	r.HandleFunc("/mobileservice", handler.List).Methods("GET")
}

//TODO maybe better place to put this
func handleCommonErrorCases(err error, rw http.ResponseWriter, logger *logrus.Logger) {
	err = errors.Cause(err)
	if mobile.IsNotFoundError(err) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	if e, ok := err.(*mobile.StatusError); ok {
		logger.Error(fmt.Sprintf("status error occurred %+v", err))
		http.Error(rw, err.Error(), e.StatusCode())
		return
	}
	if e, ok := err.(*kerror.StatusError); ok {
		logger.Error(fmt.Sprintf("kubernetes status error occurred %+v", err))
		http.Error(rw, e.Error(), int(e.Status().Code))
		return
	}
	logger.Error(fmt.Sprintf("unexpected and unknown error occurred %+v", err))
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}
