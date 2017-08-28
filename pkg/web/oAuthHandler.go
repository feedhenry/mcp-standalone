package web

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"golang.org/x/oauth2"
)

// OAuthHandler handle oauth actions
type OAuthHandler struct {
	logger *logrus.Logger
	config *oauth2.Config
}

// NewOauthHandler returns a new oauth handler
func NewOauthHandler(logger *logrus.Logger, config *oauth2.Config) *OAuthHandler {
	// OpenShift OAuth requires client id & secret in request parameters
	oauth2.RegisterBrokenAuthHeaderProvider(config.Endpoint.TokenURL)
	return &OAuthHandler{
		logger: logger,
		config: config,
	}
}

func (oah *OAuthHandler) OAuthToken(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	tr := &http.Transport{
		// TODO: skipping insecure check is OK for POC only
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	token, err := oah.config.Exchange(ctx, code)
	if err != nil {
		oah.logger.Error("Code exchange failed with ", err)
		http.Redirect(w, r, fmt.Sprintf("%s/error?error=code_exchange_failed&error_description=%s", oah.config.RedirectURL, err), http.StatusTemporaryRedirect)
		return
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJSON)
}
