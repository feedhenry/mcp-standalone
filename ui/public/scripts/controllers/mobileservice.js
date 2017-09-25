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
          if ($scope && $scope.service) {
            if ($scope.service.name !== 'fh-sync-server') {
              charts.forEach(chart => {
                c3.generate(chart);
              });
            } else {
              sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-sparkline-1';
              sparklineConfig.data = {
                columns: [
                  [
                    '%',
                    10,
                    50,
                    28,
                    20,
                    31,
                    27,
                    60,
                    36,
                    52,
                    55,
                    62,
                    68,
                    69,
                    88,
                    74,
                    88,
                    95
                  ]
                ],
                type: 'area-spline'
              };
              c3.generate(sparklineConfig);

              sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-sparkline-2';
              sparklineConfig.color = {
                pattern: ['#2ca02c']
              };
              sparklineConfig.data = {
                columns: [
                  [
                    '%',
                    10,
                    50,
                    28,
                    20,
                    31,
                    27,
                    60,
                    36,
                    52,
                    55,
                    62,
                    68,
                    69,
                    88,
                    74,
                    88,
                    95
                  ]
                ],
                type: 'area-spline'
              };
              c3.generate(sparklineConfig);

              sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-sparkline-3';
              sparklineConfig.color = {
                pattern: ['#ff7f0e']
              };
              sparklineConfig.data = {
                columns: [
                  [
                    '%',
                    10,
                    50,
                    28,
                    20,
                    31,
                    27,
                    60,
                    36,
                    52,
                    55,
                    62,
                    68,
                    69,
                    88,
                    74,
                    88,
                    95
                  ]
                ],
                type: 'area-spline'
              };
              c3.generate(sparklineConfig);

              var sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-sparkline-6';
              sparklineConfig.data = {
                columns: [
                  [
                    '%',
                    10,
                    50,
                    28,
                    20,
                    31,
                    27,
                    60,
                    36,
                    52,
                    55,
                    62,
                    68,
                    69,
                    88,
                    74,
                    88,
                    95
                  ]
                ],
                type: 'area-spline'
              };
              var chart2 = c3.generate(sparklineConfig);

              sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-sparkline-7';
              sparklineConfig.color = {
                pattern: ['#2ca02c']
              };
              sparklineConfig.data = {
                columns: [
                  [
                    '%',
                    35,
                    36,
                    20,
                    30,
                    31,
                    22,
                    44,
                    36,
                    40,
                    41,
                    55,
                    52,
                    48,
                    48,
                    50,
                    40,
                    41
                  ]
                ],
                type: 'area-spline'
              };
              var chart4 = c3.generate(sparklineConfig);

              sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-sparkline-8';
              sparklineConfig.color = {
                pattern: ['#ff7f0e']
              };
              sparklineConfig.data = {
                columns: [
                  [
                    '%',
                    60,
                    55,
                    70,
                    44,
                    31,
                    67,
                    54,
                    46,
                    58,
                    75,
                    62,
                    68,
                    69,
                    88,
                    74,
                    88,
                    85
                  ]
                ],
                type: 'area-spline'
              };
              var chart6 = c3.generate(sparklineConfig);

              charts[0].data.columns.push(charts[1].data.columns[1]);
              charts[0].data.columns.push(charts[2].data.columns[1]);
              charts[0].data.columns.push(charts[3].data.columns[1]);
              charts[0].data.columns.push(charts[4].data.columns[1]);
              charts[0].data.columns.push(charts[5].data.columns[1]);
              charts[0].data.columns.push(charts[6].data.columns[1]);
              charts[0].data.columns.push(charts[7].data.columns[1]);
              charts[0].data.columns.push(charts[8].data.columns[1]);
              charts[0].data.columns.push(charts[9].data.columns[1]);
              charts[0].data.columns.push(charts[10].data.columns[1]);
              charts[0].data.columns.push(charts[11].data.columns[1]);
              charts[0].data.type = 'spline';
              charts[0].point.r = 2;
              c3.generate(charts[0]);
            }
          }
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
  }
]);
