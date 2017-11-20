package web

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

// StaticHandler handle static routes
type StaticHandler struct {
	logger          *logrus.Logger
	staticDirectory string
	prefix          string
	fallback        string
}

// NewStaticHandler returns a new static handler
func NewStaticHandler(logger *logrus.Logger, staticDirectory string, prefix string, fallback string) *StaticHandler {
	return &StaticHandler{
		logger:          logger,
		staticDirectory: staticDirectory,
		prefix:          prefix,
		fallback:        fallback,
	}
}

func (sh StaticHandler) Static(res http.ResponseWriter, req *http.Request) {
	dir := http.Dir(sh.staticDirectory)
	if req.Method != "GET" && req.Method != "HEAD" {
		http.Error(res, "Not Found", 404)
		return
	}
	file := req.URL.Path
	// if we have a prefix, filter requests by stripping the prefix
	// /../../
	f, err := dir.Open(file)
	if err != nil {
		// try fallback before giving up
		f, err = dir.Open(sh.fallback)

		if err != nil {
			handleCommonErrorCases(err, res, sh.logger)
			return
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			sh.logger.Error("failed to close file handle. Could be leaking resources ", err)
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		handleCommonErrorCases(err, res, sh.logger)
		return
	}

	if fi.IsDir() {
		// try fallback before giving up
		f, err = dir.Open(sh.fallback)

		if err != nil {
			handleCommonErrorCases(err, res, sh.logger)
			return
		}
	}

	http.ServeContent(res, req, file, fi.ModTime(), f)
}
