package web

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/feedhenry/mobile-server/pkg/web/middleware"
	"github.com/gorilla/mux"
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
	n.UseFunc(access.Handle)
	n.UseHandler(r)
	return n
}

// MobileAppRoute configure and setup the /mobileapp route. The middleware.Builder is responsible for building per request instances of clients
func MobileAppRoute(r *mux.Router, handler *MobileAppHandler, clientBuilderMiddlware *middleware.Builder) {
	r.HandleFunc("/mobileapp", clientBuilderMiddlware.HandleRepo(handler.Create)).Methods("POST")
	r.HandleFunc("/mobileapp/{id}", clientBuilderMiddlware.HandleRepo(handler.Delete)).Methods("DELETE")
	r.HandleFunc("/mobileapp/{id}", clientBuilderMiddlware.HandleRepo(handler.Read)).Methods("GET")
	r.HandleFunc("/mobileapp", clientBuilderMiddlware.HandleRepo(handler.List)).Methods("GET")
	r.HandleFunc("/mobileapp/{id}", clientBuilderMiddlware.HandleRepo(handler.Update)).Methods("PUT")
}

// SysRoute congifures and sets up the /sys/* route
func SysRoute(r *mux.Router, handler *SysHandler) {
	r.HandleFunc("/sys/info/ping", handler.Ping).Methods("GET")
	r.HandleFunc("/sys/info/health", handler.Health).Methods("GET")
}
