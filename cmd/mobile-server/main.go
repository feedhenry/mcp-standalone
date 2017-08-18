package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/data"
	"github.com/feedhenry/mobile-server/pkg/k8s"
	"github.com/feedhenry/mobile-server/pkg/openshift"
	"github.com/feedhenry/mobile-server/pkg/web"
	"github.com/feedhenry/mobile-server/pkg/web/middleware"
)

func main() {
	var (
		router          = web.NewRouter()
		port            = flag.String("port", ":3001", "set the port to listen on")
		namespace       = flag.String("namespace", "", "the namespace to target")
		k8host          string
		logger          = logrus.New()
		appRepoBuilder  = &data.MobileAppRepoBuilder{}
		k8ClientBuilder = &k8s.ClientBuilder{}
	)
	flag.StringVar(&k8host, "k8-host", "", "kubernetes target")
	flag.Parse()
	if *namespace == "" {
		logger.Fatal("-namespace is a required flag")
	}
	k8ClientBuilder = k8ClientBuilder.WithNamespace(*namespace)

	if k8host == "" {
		k8host = os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT")
	}
	k8ClientBuilder = k8ClientBuilder.WithHost(k8host)
	var (
		mwBuilder     = middleware.NewBuilder(k8ClientBuilder, appRepoBuilder, *namespace)
		openshiftUser = openshift.UserAccess{Logger: logger}
		mwAccess      = middleware.NewAccess(logger, k8host, openshiftUser.ReadUserFromToken)
	)

	//mobileapp handler
	{
		appHandler := web.NewMobileAppHandler(logger)
		web.MobileAppRoute(router, appHandler, mwBuilder)
	}

	handler := web.BuildHTTPHandler(router, mwAccess)
	logger.Info("starting server on port " + *port)
	if err := http.ListenAndServe(*port, handler); err != nil {
		panic(err)
	}
}
