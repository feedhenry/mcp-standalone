package web

import (
	"net/http"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/web/headers"
	"github.com/gorilla/mux"
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

// GenerateKeys will parse the request and hand it off to the service logic to setup a new public private key pair
func (bh *BuildHandler) GenerateKeys(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	params := mux.Vars(req)
	buildID := params["buildID"]
	if buildID == "" {
		http.Error(rw, "buildID cannot be empty ", http.StatusBadRequest)
		return
	}
	buildRepo, err := bh.buildRepoBuilder.WithToken(token).Build()
	if err != nil {
		err = errors.Wrap(err, "build handler failed to create build repo instance")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	asset, _, err := bh.buildService.CreateBuildSrcKeySecret(buildRepo, buildID)
	if err != nil {
		err = errors.Wrap(err, "failed to generate keys")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	res := map[string]string{"name": asset}
	encoder := json.NewEncoder(rw)
	rw.WriteHeader(http.StatusCreated)
	if err := encoder.Encode(res); err != nil {
		err = errors.Wrap(err, "failed to encode response after creating source keys")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
}

func (bh *BuildHandler) AddAsset(rw http.ResponseWriter, req *http.Request) {
	token := headers.DefaultTokenRetriever(req.Header)
	buildAsset := &mobile.BuildAsset{}
	req.ParseMultipartForm(10 * 1000000) //10MB
	params := mux.Vars(req)
	file, info, err := req.FormFile("asset")
	if err != nil {
		err = errors.Wrap(err, "getting the form file failed when adding asset")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	defer file.Close()
	buildAsset.Name = info.Filename
	buildAsset.Platform = params["platform"]
	buildAsset.Type = mobile.BuildAssetTypeBuildSecret
	if err := buildAsset.Validate(mobile.BuildAssetTypeBuildSecret); err != nil {
		err = &mobile.StatusError{Message: err.Error(), Code: http.StatusBadRequest}
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	br, err := bh.buildRepoBuilder.WithToken(token).Build()
	if err != nil {
		err = errors.Wrap(err, "failed to create build repo with token")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	asset, err := bh.buildService.AddBuildAsset(br, file, buildAsset)
	if err != nil {
		err = errors.Wrap(err, "AddAsset failed to add new build resource")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
	res := map[string]string{"name": asset}
	encoder := json.NewEncoder(rw)
	rw.WriteHeader(http.StatusCreated)
	if err := encoder.Encode(res); err != nil {
		err = errors.Wrap(err, "failed to encode response after creating build asset")
		handleCommonErrorCases(err, rw, bh.logger)
		return
	}
}
