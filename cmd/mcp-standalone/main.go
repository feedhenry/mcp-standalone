package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/k8s"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/client"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/openshift"
	"github.com/feedhenry/mcp-standalone/pkg/web"
	"github.com/feedhenry/mcp-standalone/pkg/web/middleware"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func main() {
	var (
		router          = web.NewRouter()
		port            = flag.String("port", ":3001", "set the port to listen on")
		cert            = flag.String("cert", "server.crt", "SSL/TLS Certificate to HTTPS")
		key             = flag.String("key", "server.key", "SSL/TLS Private Key for the Certificate")
		namespace       = flag.String("namespace", os.Getenv("NAMESPACE"), "the namespace to target")
		saTokenPath     = flag.String("satoken-path", "var/run/secrets/kubernetes.io/serviceaccount/token", "where on disk the service account token to use is ")
		staticDirectory = flag.String("web-dir", "./web/app", "Location of static content to serve at /console. index.html will be used as a fallback for requested files that don't exist")
		k8host          string
		logger          = logrus.New()
		appRepoBuilder  = &data.MobileAppRepoBuilder{}
		svcRepoBuilder  = &data.MobileServiceRepoBuilder{}
	)
	flag.StringVar(&k8host, "k8-host", "", "kubernetes target")
	flag.Parse()

	if *namespace == "" {
		logger.Fatal("-namespace is a required flag or it can be set via NAMESPACE env var")
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
		httpClientBuilder  = client.NewHttpClientBuilder()
		openshiftUser      = openshift.UserAccess{Logger: logger}
		mwAccess           = middleware.NewAccess(logger, k8host, openshiftUser.ReadUserFromToken)
	)
	tokenClientBuilder.SAToken = token

	//oauth handler
	{
		kubernetesOauthEndpoint := &oauth2.Endpoint{
			AuthURL:  k8host + "/oauth/authorize",
			TokenURL: k8host + "/oauth/token",
		}

		kubernetesOauthConfig := &oauth2.Config{
			// TODO: how to dynamically configure this url from the Route
			RedirectURL:  "https://127.0.0.1:3001/console/oauth",
			ClientID:     fmt.Sprintf("system:serviceaccount:%s:mcp-standalone", *namespace),
			ClientSecret: token,
			Scopes:       []string{"user:info user:check-access"},
			Endpoint:     *kubernetesOauthEndpoint,
		}
		oauthHandler := web.NewOauthHandler(logger, kubernetesOauthConfig)
		web.OAuthRoute(router, oauthHandler)
	}

	//mobileapp handler
	{
		appHandler := web.NewMobileAppHandler(logger, tokenClientBuilder)
		web.MobileAppRoute(router, appHandler)
	}

	//mobileservice handler
	{
		integrationSvc := &integration.MobileService{}
		svcHandler := web.NewMobileServiceHandler(logger, integrationSvc, tokenClientBuilder)
		web.MobileServiceRoute(router, svcHandler)
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

	//static handler
	{
		staticHandler := web.NewStaticHandler(logger, *staticDirectory, "/console", "index.html")
		web.StaticRoute(staticHandler)
	}

	//add in the rolebinding mw
	mrb := middleware.NewRoleBinding(httpClientBuilder, *namespace, logger, k8host)

	handler := web.BuildHTTPHandler(router, mwAccess, mrb)
	http.Handle("/", handler)
	logger.Info("starting server on port "+*port, " using key ", *key, " and cert ", *cert, "target namespace is ", *namespace)

	if err := http.ListenAndServeTLS(*port, *cert, *key, nil); err != nil {
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
