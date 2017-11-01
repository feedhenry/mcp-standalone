'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileServiceController
 * @description
 * # MobileServiceController
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp').controller('MobileServiceController', [
  '$scope',
  '$timeout',
  'mcpApi',
  '$routeParams',
  'ProjectsService',
  'DataService',
  function(
    $scope,
    $timeout,
    mcpApi,
    $routeParams,
    ProjectsService,
    DataService
  ) {
    $scope.alerts = {};
    $scope.projectName = $routeParams.project;
    $scope.breadcrumbs = [
      {
        title: 'Mobile Service',
        link: 'project/' + $routeParams.project + '/browse/mobileoverview'
      },
      {
        title: $routeParams.service
      }
    ];

    const knownServices = ['fh-sync-server', 'keycloak'];
    $scope.chartData = [];
    $scope.integrations = [];
    mcpApi
      .mobileServiceMetrics($routeParams.service)
      .then(chartData => {
        $scope.chartData = chartData;
        $timeout(() => {});
      })
      .catch(err => {
        console.error('Error loading Service Metrics', err);
      });
    ProjectsService.get($routeParams.project).then(projectInfo => {
      const [project = {}, projectContext = {}] = projectInfo;
      $scope.project = project;
      $scope.projectContext = projectContext;
      DataService.list(
        {
          group: 'servicecatalog.k8s.io',
          resource: 'clusterserviceclasses'
        },
        $scope.projectContext,
        function(serviceClasses) {
          $scope.serviceClasses = serviceClasses['_data'];
        }
      );
    });

    Promise.all([
      mcpApi.mobileService($routeParams.service, 'true').then(s => {
        $scope.service = s;
        $scope.integrations = Object.keys(s.integrations);
      }),
      mcpApi.mobileApps().then(apps => {
        $scope.mobileappsCount = apps.length;
        $scope.mobileapps = {};
        for (var i = 0; i < apps.length; i++) {
          let app = apps[i];
          $scope.mobileapps[app.clientType] = 'true';
        }
        $scope.clients = Object.keys($scope.mobileapps);
        $scope.clientType = $scope.clients[0];
      })
    ]).then(() => {
      // wait for apps & services, and hence the UI being redrawn,
      // before fetching metrics.
      // Charts may not initialise if the UI isn't ready with the div placeholders
      mcpApi.mobileServiceMetrics($routeParams.service).then(data => {
        var charts = [];
        var chartConfigs = [];
        data.forEach(columns => {
          var c3ChartDefaults = $().c3ChartDefaults();
          var chartConfig = c3ChartDefaults.getDefaultLineConfig();
          chartConfig.axis = {
            x: {
              type: 'timeseries',
              tick: {
                // 11:34:55
                format: '%H:%M:%S'
              }
            }
          };
          chartConfig.data = {
            x: 'x',
            xFormat: '%Y-%m-%d %H:%M:%S',
            columns: columns,
            type: 'line'
          };
          var chartName = columns[1][0];
          chartConfig.bindto = '#line-chart-' + chartName;
          chartConfig.title = chartName;
          charts.push(chartConfig);
        });

        $scope.charts = charts;
      });
    });

    Promise.all([
      mcpApi.mobileService($routeParams.service, 'true'),
      mcpApi.mobileApps()
    ])
      .then(serviceInfo => {
        const [service = {}, apps = []] = serviceInfo;

        $scope.service = service;
        $scope.integrations = Object.keys(service.integrations);

        $scope.templateName = knownServices.includes(service.name)
          ? service.name
          : 'default-service';

        $scope.mobileappsCount = apps.length;
        $scope.mobileapps = apps.reduce((acc, current) => {
          acc[current.clientType] = 'true';
          return acc;
        }, {});
        $scope.clients = Object.keys($scope.mobileapps);
        $scope.clientType = $scope.clients[0];
        $timeout(() => {});
      })
      .catch(err => {
        console.error('Error loading Services', err);
      });

    $scope.status = function(integration, service) {
      if ($scope.processing(integration, service)) {
        return 2;
      }
      if ($scope.enabled(integration, service)) {
        return 1;
      }
      return 0;
    };
    $scope.serviceWritable = function(service) {
      return service.writable;
    };
    $scope.enabled = function(integration, service) {
      if (!service) {
        return false;
      }
      if (service.integrations[integration]) {
        return service.integrations[integration].enabled === true;
      }
      return false;
    };
    $scope.processing = function(integration, service) {
      if (!service) {
        return false;
      }
      if (
        service.integrations[integration] &&
        service.integrations[integration].processing
      ) {
        return service.integrations[integration].processing === true;
      }
      return false;
    };
    $scope.enableIntegration = function(service, key) {
      var integration = service.integrations[key];
      integration.processing = true;
      mcpApi
        .integrateService(service, key)
        .then(res => {
          integration.processing = false;
          integration.enabled = true;
        })
        .catch(e => {
          integration.processing = false;
          console.log('error integrating service ', e);
        });
      return true;
    };
    $scope.disableIntegration = function(service, key) {
      var integration = service.integrations[key];
      service.processing = true;
      mcpApi
        .deintegrateService(service, key)
        .then(res => {
          integration.processing = false;
          integration.enabled = false;
        })
        .catch(e => {
          integration.processing = false;
          console.log('error deintegrating service ', e);
        });
      return true;
    };

    $scope.clientType =
      $scope.clients && $scope.clients.length > 0
        ? $scope.clients[0]
        : 'cordova';

    $scope.installationOpt = function(type) {
      $scope.clientType = type;
    };

    $scope.dashboardSelected = function() {
      window.dispatchEvent(new Event('resize'));
    };

    $scope.getDescription = function(service) {
      for (var scId in $scope.serviceClasses) {
        var sc = $scope.serviceClasses[scId];
        if (
          sc.spec.externalMetadata.hasOwnProperty('serviceName') &&
          sc.spec.externalMetadata.serviceName.toLowerCase() ==
            service.type.toLowerCase()
        ) {
          return sc.spec.description;
        }
      }
      return '';
    };
  }
]);
