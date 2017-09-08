package web

import (
	"net/http"
	"text/template"

	"bytes"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
)

// ConsoleConfigHandler handle console config route
type ConsoleConfigHandler struct {
	logger           *logrus.Logger
	consoleMountPath string
	mcpConsoleConfig mcpConsoleConfig
}

var configTemplate = template.Must(template.New("mcpConsoleConfig").Parse(`
window.OPENSHIFT_CONFIG = {
  apis: {
    hostPort: "{{ .APIGroupAddr | js}}",
    prefix: "{{ .APIGroupPrefix | js}}"
  },
  api: {
    openshift: {
      hostPort: "{{ .MasterAddr | js}}",
      prefix: "{{ .MasterPrefix | js}}"
    },
    k8s: {
      hostPort: "{{ .KubernetesAddr | js}}",
      prefix: "{{ .KubernetesPrefix | js}}"
    }
  },
auth: {
  	oauth_authorize_uri: "{{ .OAuthAuthorizeURI | js}}",
	oauth_token_uri: "{{ .OAuthTokenURI | js}}",
  	oauth_redirect_base: "{{ .OAuthRedirectBase | js}}",
  	oauth_client_id: "{{ .OAuthClientID | js}}",
  	logout_uri: "{{ .LogoutURI | js}}",
  	scope: "{{ .Scope | js}}"
}
};
`))

type mcpConsoleConfig struct {
	// APIGroupAddr is the host:port the UI should call the API groups on. Scheme is derived from the scheme the UI is served on, so they must be the same.
	APIGroupAddr string
	// APIGroupPrefix is the API group context root
	APIGroupPrefix string
	// MasterAddr is the host:port the UI should call the master API on. Scheme is derived from the scheme the UI is served on, so they must be the same.
	MasterAddr string
	// MasterPrefix is the OpenShift API context root
	MasterPrefix string
	// KubernetesAddr is the host:port the UI should call the kubernetes API on. Scheme is derived from the scheme the UI is served on, so they must be the same.
	// TODO this is probably unneeded since everything goes through the openshift master's proxy
	KubernetesAddr string
	// KubernetesPrefix is the Kubernetes API context root
	KubernetesPrefix string
	// OAuthAuthorizeURI is the OAuth2 endpoint to use to request an API token. It must support request_type=token.
	OAuthAuthorizeURI string
	// OAuthTokenURI is the OAuth2 endpoint to use to request an API token. If set, the OAuthClientID must support a client_secret of "".
	OAuthTokenURI string
	// OAuthRedirectBase is the base URI of the web console. It must be a valid redirect_uri for the OAuthClientID
	OAuthRedirectBase string
	// OAuthClientID is the OAuth2 client_id to use to request an API token. It must be authorized to redirect to the web console URL.
	OAuthClientID string
	// LogoutURI is an optional (absolute) URI to redirect to after completing a logout. If not specified, the built-in logout page is shown.
	LogoutURI string
	// Oauth Scopes to request a code/token for
	Scope string
}

// NewConsoleConfigHandler returns a new console config handler
func NewConsoleConfigHandler(logger *logrus.Logger, k8sHost, k8sAuthorizeEndpoint, oauthClientID, namespace string) *ConsoleConfigHandler {
	mcpConsoleConfig := mcpConsoleConfig{
		APIGroupAddr:      k8sHost,
		APIGroupPrefix:    "/apis",
		MasterAddr:        k8sHost,
		MasterPrefix:      "/oapi",
		KubernetesAddr:    k8sHost,
		KubernetesPrefix:  "/api",
		OAuthAuthorizeURI: k8sAuthorizeEndpoint,
		OAuthTokenURI:     "",
		OAuthRedirectBase: "",
		OAuthClientID:     oauthClientID,
		LogoutURI:         "",
		Scope:             fmt.Sprintf("user:info user:check-access role:edit:%s:!", namespace),
	}

	return &ConsoleConfigHandler{
		logger:           logger,
		mcpConsoleConfig: mcpConsoleConfig,
	}
}

func (cch ConsoleConfigHandler) Config(res http.ResponseWriter, req *http.Request) {
	baseUrl, err := headers.ParseBaseUrl(req)
	if err != nil {
		handleCommonErrorCases(err, res, cch.logger)
	}
	cch.mcpConsoleConfig.OAuthTokenURI = fmt.Sprintf("%s/oauth/token", baseUrl)
	cch.mcpConsoleConfig.OAuthRedirectBase = baseUrl

	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, cch.mcpConsoleConfig); err != nil {
		handleCommonErrorCases(errors.Wrap(err, "Error executing console config template"), res, cch.logger)
	}

	config := buffer.Bytes()

	res.Header().Add("Cache-Control", "no-cache, no-store")
	res.Header().Add("Content-Type", "application/javascript")
	if _, err := res.Write(config); err != nil {
		handleCommonErrorCases(err, res, cch.logger)
		return
	}
}
