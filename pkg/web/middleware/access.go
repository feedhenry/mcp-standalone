package middleware

import (
	"fmt"
	"net/http"
	"regexp"

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

func buildIgnoreList() []*regexp.Regexp {
	r := regexp.MustCompile("^/sdk/mobileapp/.*/config")
	return []*regexp.Regexp{
		r,
	}
}

var ingnoreList = buildIgnoreList()

func (c Access) shouldIgnore(path string) bool {
	for _, i := range ingnoreList {
		if i.Match([]byte(path)) {
			c.logger.Info("ignoring user access check on path: ", path)
			return true
		}
	}
	return false
}

// Handle sets the required headers
func (c Access) Handle(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token := req.Header.Get("x-auth")
	if c.shouldIgnore(req.URL.Path) {
		next(w, req)
		return
	}
	//todo take config to set skipTLS
	if err := c.userCheck(c.host, token, true); err != nil {
		if openshift.IsAuthenticationError(err) {
			c.logger.Error(err.Error(), err)
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
