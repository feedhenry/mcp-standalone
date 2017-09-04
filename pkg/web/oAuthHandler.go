package web

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/k8s"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"golang.org/x/oauth2"
)

// OAuthHandler handle oauth actions
type OAuthHandler struct {
	logger        *logrus.Logger
	k8sMetadata   k8s.Metadata
	oauthClientID string
	oauthEndpoint oauth2.Endpoint
	token         string
}

// NewOauthHandler returns a new oauth handler
func NewOauthHandler(logger *logrus.Logger, k8sMetadata k8s.Metadata, oauthClientID string, token string) *OAuthHandler {
	// OpenShift OAuth requires client id & secret in request parameters
	oauthEndpoint := oauth2.Endpoint{
		AuthURL:  k8sMetadata.AuthorizationEndpoint,
		TokenURL: k8sMetadata.TokenEndpoint,
	}
	oauth2.RegisterBrokenAuthHeaderProvider(oauthEndpoint.TokenURL)

	return &OAuthHandler{
		logger:        logger,
		k8sMetadata:   k8sMetadata,
		oauthClientID: oauthClientID,
		oauthEndpoint: oauthEndpoint,
		token:         token,
	}
}

func (oah *OAuthHandler) OAuthToken(w http.ResponseWriter, r *http.Request) {
	baseUrl, err := headers.ParseBaseUrl(r)
	if err != nil {
		handleCommonErrorCases(err, w, oah.logger)
	}

	config := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/console/oauth", baseUrl),
		ClientID:     oah.oauthClientID,
		ClientSecret: oah.token,
		Scopes:       []string{"user:info user:check-access"},
		Endpoint:     oah.oauthEndpoint,
	}

	code := r.FormValue("code")

	tr := &http.Transport{
		// TODO: skipping insecure check is OK for POC only
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		oah.logger.Error("Code exchange failed with ", err)
		http.Redirect(w, r, fmt.Sprintf("%s/error?error=code_exchange_failed&error_description=%s", config.RedirectURL, err), http.StatusTemporaryRedirect)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(token); err != nil {
		handleCommonErrorCases(err, w, oah.logger)
		return
	}
}
