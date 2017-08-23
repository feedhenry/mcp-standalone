package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mobile-server/pkg/data"
	"github.com/feedhenry/mobile-server/pkg/k8s"
	"github.com/feedhenry/mobile-server/pkg/mobile/client"
	"github.com/feedhenry/mobile-server/pkg/mobile/integration"
	"github.com/feedhenry/mobile-server/pkg/openshift"
	"github.com/feedhenry/mobile-server/pkg/web"
	"github.com/feedhenry/mobile-server/pkg/web/middleware"
	"github.com/pkg/errors"
)

func main() {
	var (
		router         = web.NewRouter()
		port           = flag.String("port", ":3001", "set the port to listen on")
		namespace      = flag.String("namespace", "", "the namespace to target")
		saTokenPath    = flag.String("satoken-path", "var/run/secrets/kubernetes.io/serviceaccount/token", "where on disk the service account token to use is ")
		k8host         string
		logger         = logrus.New()
		appRepoBuilder = &data.MobileAppRepoBuilder{}
		svcRepoBuilder = &data.MobileServiceRepoBuilder{}
	)
	flag.StringVar(&k8host, "k8-host", "", "kubernetes target")
	flag.Parse()
	if *namespace == "" {
		logger.Fatal("-namespace is a required flag")
	}

	token, err := readSAToken(*saTokenPath)
	if err != nil {
		panic(err)
	}

	if k8host == "" {
		k8host = "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT")
	}
	var k8ClientBuilder = k8s.NewClientBuilder(*namespace, k8host)
	var (
		tokenClientBuilder = client.NewTokenScopedClientBuilder(k8ClientBuilder, appRepoBuilder, svcRepoBuilder, *namespace, logger)
		openshiftUser      = openshift.UserAccess{Logger: logger}
		mwAccess           = middleware.NewAccess(logger, k8host, openshiftUser.ReadUserFromToken)
	)
	tokenClientBuilder.SAToken = token

	//mobileapp handler
	{
		appHandler := web.NewMobileAppHandler(logger, tokenClientBuilder)
		web.MobileAppRoute(router, appHandler)
	}
	//sdk handler
	{
		integrationSvc := &integration.MobileService{}
		sdkHandler := web.NewSDKConfigHandler(logger, integrationSvc, tokenClientBuilder)
		web.SDKConfigRoute(router, sdkHandler)
	}
	//sys handler
	{
		sysHandler := web.NewSysHandler(logger)
		web.SysRoute(router, sysHandler)
	}

	handler := web.BuildHTTPHandler(router, mwAccess)
	logger.Info("starting server on port " + *port)
	if err := http.ListenAndServe(*port, handler); err != nil {
		panic(err)
	}
}

func readSAToken(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to read service account token ")
	}
	return string(data), nil
}
