package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/web"
)

func setupSysHandler() http.Handler {
	router := web.NewRouter()
	httpHandler := web.BuildHTTPHandler(router)
	sysHandler := web.NewSysHandler(logrus.StandardLogger())
	web.SysRoute(router, sysHandler)
	return httpHandler
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(setupSysHandler())
	defer server.Close()
}
