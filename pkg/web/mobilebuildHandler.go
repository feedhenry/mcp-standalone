package web

import (
	"net/http"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/pkg/errors"
)

type BuildHandler struct {
	buildRepoBuilder mobile.BuildRepoBuilder
	buildService     *app.Build
	logger           *logrus.Logger
}

// NewBuildHandler returns a configured build handler
func NewBuildHandler(br mobile.BuildRepoBuilder, buildService *app.Build, logger *logrus.Logger) *BuildHandler {
	return &BuildHandler{
		buildRepoBuilder: br,
		buildService:     buildService,
		logger:           logger,
	}
}

// Create will parse the create request and hand it off to app build service
func (bh *BuildHandler) Create(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	buildRepo, err := bh.buildRepoBuilder.WithToken(token).Build()
	if err != nil {
		err = errors.Wrap(err, "build handler failed to create build repo instance")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	var (
		build   = &mobile.Build{}
		decoder = json.NewDecoder(req.Body)
		encoder = json.NewEncoder(rw)
	)

	if err := decoder.Decode(build); err != nil {
		err = errors.Wrap(err, "build handler failed to decode build payload")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}

	res, err := bh.buildService.CreateAppBuild(buildRepo, build)
	if err != nil {
		err = errors.Wrap(err, "build handler failed to create app build")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	if err := encoder.Encode(res); err != nil {
		err = errors.Wrap(err, "failed to encode the build response")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}

}
