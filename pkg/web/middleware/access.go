package middleware

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/openshift"
)

type UserChecker func(host, token string, skipTLS bool) error

// Access handles the cross origin requests
type Access struct {
	host      string
	logger    *logrus.Logger
	userCheck UserChecker
}

func NewAccess(logger *logrus.Logger, host string, userCheck UserChecker) *Access {
	return &Access{
		logger:    logger,
		host:      host,
		userCheck: userCheck,
	}
}

// Handle sets the required headers
func (c Access) Handle(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token := req.Header.Get("x-auth")
	//todo take config to set skipTLS
	if err := c.userCheck(c.host, token, true); err != nil {
		if openshift.IsAuthenticationError(err) {
			c.logger.Error(fmt.Sprintf("access was denied %s : %+v", err.Error(), err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		c.logger.Error(fmt.Sprintf("failed to read user to check access %s : %+v", err.Error(), err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, req)
	return

}
