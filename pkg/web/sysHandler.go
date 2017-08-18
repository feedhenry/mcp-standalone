package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// SysHandler handles the /sys route
type SysHandler struct {
	log *logrus.Logger
}

// Ping just returns a 200 ok
func (sh *SysHandler) Ping(rw http.ResponseWriter, req *http.Request) {}

// Health tells us if parts of the system are ok
func (sh *SysHandler) Health(rw http.ResponseWriter, req *http.Request) {
	healthChecks := map[string]string{
		"http": "ok",
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(healthChecks); err != nil {
		http.Error(rw, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewSysHandler(logger *logrus.Logger) *SysHandler {
	return &SysHandler{
		log: logger,
	}
}
