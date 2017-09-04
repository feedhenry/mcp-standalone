package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/pkg/errors"
)

// MobileServiceHandler handles endpoints associated with mobile enabled services. It will list services in the namespace that
// it knows about so that they can be rendered in the MCP
type MobileServiceHandler struct {
	logger                   *logrus.Logger
	mobileIntegrationService *integration.MobileService
	tokenClientBuilder       mobile.TokenScopedClientBuilder
}

// NewMobileServiceHandler returns a new MobileServiceHandler
func NewMobileServiceHandler(logger *logrus.Logger, integrationService *integration.MobileService, tokenClientBuilder mobile.TokenScopedClientBuilder) *MobileServiceHandler {
	return &MobileServiceHandler{
		logger: logger,
		mobileIntegrationService: integrationService,
		tokenClientBuilder:       tokenClientBuilder,
	}
}

// List allows you to list mobile services
func (msh *MobileServiceHandler) List(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	serviceCruder, err := msh.tokenClientBuilder.MobileServiceCruder(token)
	if err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	svc, err := msh.mobileIntegrationService.DiscoverMobileServices(serviceCruder)
	if err != nil {
		err = errors.Wrap(err, "attempted to list mobile services")
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(svc); err != nil {
		handleCommonErrorCases(err, rw, msh.logger)
		return
	}
}
