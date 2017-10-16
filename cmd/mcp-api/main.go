package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"syscall"

	"context"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/data"
	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/feedhenry/mcp-standalone/pkg/jenkins"
	"github.com/feedhenry/mcp-standalone/pkg/k8s"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/app"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/integration"
	"github.com/feedhenry/mcp-standalone/pkg/mobile/metrics"
	"github.com/feedhenry/mcp-standalone/pkg/openshift"
	"github.com/feedhenry/mcp-standalone/pkg/web"
	"github.com/feedhenry/mcp-standalone/pkg/web/middleware"
	"github.com/pkg/errors"
)

var (
	// Version is set at build time using ldflags see the Makefile for reference
	Version string
)

func main() {
	var (
		k8host             string
		port               = flag.String("port", ":3001", "set the port to listen on")
		insecure           = flag.String("insecure", "false", "allow insecure requests")
		cert               = flag.String("cert", "server.crt", "SSL/TLS Certificate to HTTPS")
		key                = flag.String("key", "server.key", "SSL/TLS Private Key for the Certificate")
		namespace          = flag.String("namespace", os.Getenv("NAMESPACE"), "the namespace to target")
		logLevel           = flag.String("log-level", "error", "the level to log at")
		saTokenPath        = flag.String("satoken-path", "var/run/secrets/kubernetes.io/serviceaccount/token", "where on disk the service account token to use is ")
		staticDirectory    = flag.String("web-dir", "./web/app", "Location of static content to serve at /console. index.html will be used as a fallback for requested files that don't exist")
		serverReadTimeout  = flag.Int("server-read-timeout", 5, "sets the maximum amount of time in seconds allowed for the server to read the headers")
		serverWriteTimeout = flag.Int("server-write-timeout", 15, " WriteTimeout is the maximum duration in seconds before timing out writes of the response")
		serverIdleTimeout  = flag.Int("server-idle-timeout", 60, " IdleTimeout is the maximum amount of time in seconds to wait for the next request when keep-alives are enabled. ")
	)
	flag.StringVar(&k8host, "k8-host", "", "kubernetes target")
	flag.Parse()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
	logger := logrus.StandardLogger()

	if *namespace == "" {
		logger.Fatal("-namespace is a required flag or it can be set via NAMESPACE env var")
	}

	token, err := readSAToken(*saTokenPath)
	if err != nil {
		logger.Fatal("reading the service account token failed ", err)
	}

	if k8host == "" {
		k8host = "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT")
	}

	logger.Info("insecure requests set to ", *insecure)

	var (
		router           = web.NewRouter()
		insecureRequests = *insecure == "true"
		incluster        = os.Getenv("KUBERNETES_SERVICE_HOST") != ""
		//setup our builders
		k8ClientBuilder    = k8s.NewClientBuilder(*namespace, k8host, insecureRequests)
		appRepoBuilder     = data.NewMobileAppRepoBuilder(k8ClientBuilder, *namespace, token)
		svcRepoBuilder     = data.NewServiceRepoBuilder(k8ClientBuilder, *namespace, token)
		authCheckerBuilder = openshift.NewAuthCheckerBuilder(k8host)
		userRepoBuilder    = openshift.NewUserRepoBuilder(k8host, insecureRequests).WithClient(&openshift.UserAccess{})
		httpClientBuilder  = httpclient.NewClientBuilder()
		defaultHTTPClient  = httpClientBuilder.Insecure(insecureRequests).Build()
		ocClientBuilder    = openshift.NewClientBuilder(k8host, *namespace, incluster, insecureRequests)
		buildRepoBuilder   = data.NewBuildsRepoBuilder(k8ClientBuilder, ocClientBuilder, *namespace, token)
		openshiftUser      = openshift.UserAccess{}
		mwAccess           = middleware.NewAccess(logger, k8host, openshiftUser.ReadUserFromToken)
		// these channels control when background proccess should stop
		stop = make(chan struct{})
		s    = make(chan os.Signal, 1)
	)

	// send a message to the signal channel for any interrupt type signals (ctl+c etc)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	k8sMetadata, err := k8s.GetMetadata(k8host, defaultHTTPClient)
	if err != nil {
		logger.Fatal("failed to get Kubernetes meta data ", err)
	}
	// Ensure that the apiKey map exists
	{
		err := createAppAPIKeyMap(appRepoBuilder, token)
		if err != nil {
			logger.Fatal("failed to create the appApiKeyMap at start up ", err)
		}
	}
	//kick off metrics scheduler
	{
		//TODO move time interval to config
		interval := time.NewTicker(30 * time.Second)
		gatherer := metrics.NewGathererScheduler(interval, stop, logger)

		// add metrics gatherers
		kcMetrics := metrics.NewKeycloak(httpClientBuilder, svcRepoBuilder, logger)
		gatherer.Add(kcMetrics.ServiceName, kcMetrics.Gather)

		// add fh-sync-server gatherers
		syncMetrics := metrics.NewFhSyncServer(httpClientBuilder, svcRepoBuilder, logger)
		gatherer.Add(syncMetrics.ServiceName, syncMetrics.Gather)

		// start collecting metrics
		go gatherer.Run()
	}
	//mobileapp httpHandler
	{
		appService := &app.Service{}
		appHandler := web.NewMobileAppHandler(logger, appRepoBuilder, appService)
		web.MobileAppRoute(router, appHandler)
	}

	//mobileservice httpHandler
	{
		integrationSvc := integration.NewMobileSevice(*namespace)
		metricSvc := &metrics.MetricsService{}
		scClientBuilder := k8s.NewServiceCatalogClientBuilder(k8ClientBuilder, defaultHTTPClient, token, *namespace, k8host)
		svcHandler := web.NewMobileServiceHandler(logger, integrationSvc, metricSvc, svcRepoBuilder, userRepoBuilder, authCheckerBuilder, scClientBuilder)
		web.MobileServiceRoute(router, svcHandler)
	}
	//sdk httpHandler
	{
		sdkService := &integration.SDKService{}
		sdkHandler := web.NewSDKConfigHandler(logger, sdkService, svcRepoBuilder, appRepoBuilder)
		web.SDKConfigRoute(router, sdkHandler)
	}
	//sys httpHandler
	{
		sysHandler := web.NewSysHandler(logger)
		web.SysRoute(router, sysHandler)
	}
	//build httpHandler
	{
		artifactRet := jenkins.NewClient(defaultHTTPClient, logger)
		buildSvc := app.NewBuild(artifactRet, token)
		buildHandler := web.NewBuildHandler(buildRepoBuilder, buildSvc, logger)
		web.MobileBuildRoute(router, buildHandler)
	}

	//console config httpHandler
	var consoleMountPath = ""
	{
		k8MetaHost, err := k8sMetadata.GetK8IssuerHost()
		if err != nil {
			logger.Fatal("failed to get the k8 meta host at startup", err)
		}
		consoleConfigHandler := web.NewConsoleConfigHandler(logger, k8MetaHost, k8sMetadata.AuthorizationEndpoint, *namespace)
		web.ConsoleConfigRoute(router, consoleConfigHandler)
	}
	//static httpHandler
	{
		staticHandler := web.NewStaticHandler(logger, *staticDirectory, consoleMountPath, "index.html")
		web.StaticRoute(staticHandler)
	}
	httpHandler := web.BuildHTTPHandler(router, mwAccess)
	server := http.Server{
		Addr:              *port,
		IdleTimeout:       time.Second * time.Duration(*serverIdleTimeout),
		ReadHeaderTimeout: time.Second * time.Duration(*serverReadTimeout),
		WriteTimeout:      time.Second * time.Duration(*serverWriteTimeout),
		Handler:           httpHandler,
	}

	logger.Info("starting server on port "+*port, " using key ", *key, " and cert ", *cert, "target namespace is ", *namespace, ". Version is: "+Version)
	go func() {
		if err := server.ListenAndServeTLS(*cert, *key); err != nil {
			logger.Fatal("failed to listen and server https ", err)
		}
	}()
	<-s //wait for interrupt
	close(stop)
	logger.Println("\nShutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
}

func readSAToken(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to read service account token ")
	}
	return string(data), nil
}

func createAppAPIKeyMap(appRepoBuilder mobile.AppRepoBuilder, token string) error {
	appRepo, err := appRepoBuilder.WithToken(token).Build()
	if err != nil {
		return err
	}
	err = appRepo.CreateAPIKeyMap()
	if err != nil {
		return err
	}
	return nil
}
