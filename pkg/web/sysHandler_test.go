package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/web"
)

func setupSysHandler() http.Handler {
	router := web.NewRouter()
	httpHandler := web.BuildHTTPHandler(router, nil, nil)
	sysHandler := web.NewSysHandler(logrus.StandardLogger())
	web.SysRoute(router, sysHandler)
	return httpHandler
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(setupSysHandler())
	defer server.Close()
	resp, err := http.Get(server.URL + "/sys/info/ping")
	if err != nil {
		t.Fatalf("failed to make sys info ping request")
	}
	if resp.StatusCode != 200 {
		t.Fatal("expected a 200 response code from sys/info/ping")
	}
}
