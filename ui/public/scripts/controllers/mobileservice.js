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
  function($scope, $timeout, mcpApi, $routeParams) {
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

    $scope.charts = [];
    $scope.integrations = [];
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
    ])
      .then(() => {
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
      })
      .catch(error => {
        console.error(error);
      });

    $scope.$watch('charts', charts => {
      $timeout(
        () => {
          charts.forEach(chart => {
            c3.generate(chart);
          });
        },
        0,
        false
      );
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
    $scope.serviceWriteable = function(service) {
      return service.writeable;
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
    $scope.debug = function(service) {
      console.log(service);
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
    $scope.enableIntegration = function(service) {
      service.processing = true;
      mcpApi
        .integrateService(service)
        .then(res => {
          service.processing = false;
          service.enabled = true;
        })
        .catch(e => {
          service.processing = false;
          console.log('error integrating service ', e);
        });
      return true;
    };
    $scope.disableIntegration = function(service) {
      service.processing = true;
      mcpApi
        .deintegrateService(service)
        .then(res => {
          service.processing = false;
          service.enabled = false;
        })
        .catch(e => {
          service.processing = false;
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
  }
]);
