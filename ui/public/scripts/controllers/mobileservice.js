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
      if ($scope && $scope.service) {
        if ($scope.service.name === 'fh-sync-server') {
          // custom layout for fh-sync-server
          // render all queue & worker sparklines
          [
            ['sync_worker_queue_count', '#ff7f0e', 'total'],
            ['ack_worker_queue_count', '#2ca02c', 'total'],
            ['pending_worker_queue_count', '#1f77b4', 'total'],
            ['sync_worker_process_time_ms', '#ff7f0e', 'avg'],
            ['ack_worker_process_time_ms', '#2ca02c', 'avg'],
            ['pending_worker_process_time_ms', '#1f77b4', 'avg']
          ].forEach(metric => {
            var chart = _.findWhere(charts, {
              title: metric[0]
            });
            if (chart) {
              var sparklineConfig = $()
                .c3ChartDefaults()
                .getDefaultSparklineConfig();
              sparklineConfig.bindto = '#chart-pf-' + metric[0];
              chart.data.columns[1][0] = ''; // no suffix on hover
              sparklineConfig.color = {
                pattern: [metric[1]]
              };
              sparklineConfig.data = {
                columns: [chart.data.columns[1]],
                type: 'area'
              };
              sparklineConfig.point.r = 0;
              c3.generate(sparklineConfig);

              // set the total/avg value
              var supplementId = metric[0] + '_' + metric[2];
              var supplementChart = _.findWhere(charts, {
                title: supplementId
              });
              if (supplementChart) {
                $scope[supplementId] =
                  supplementChart.data.columns[1][
                    supplementChart.data.columns[1].length - 1
                  ];
              }
            }
          });

          // Gather all 'timings' data to render into 1 chart
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
            columns: [charts[0].data.columns[0]], // copy over x values
            type: 'spline'
          };
          chartConfig.point.r = 2;
          chartConfig.bindto = '#line-chart-timings';

          charts.forEach(chart => {
            if (
              chart.title.indexOf('mongodb_operation_time_') === 0 ||
              chart.title.indexOf('api_process_time_') === 0
            ) {
              chartConfig.data.columns.push(chart.data.columns[1]);
            }
          });

          c3.generate(chartConfig);
        } else if ($scope.service.name === 'keycloak') {
          var keycloakMetrics = {
            logins: {
              success: 0,
              error: 0
            },
            registrations: {
              success: 0,
              error: 0
            }
          };
          // custom layout for keycloak
          var loginSuccessChart = _.findWhere(charts, {
            title: 'LOGIN'
          });
          if (loginSuccessChart) {
            keycloakMetrics.logins.success =
              loginSuccessChart.data.columns[1][
                loginSuccessChart.data.columns[1].length - 1
              ];
          }
          var loginErrorChart = _.findWhere(charts, {
            title: 'LOGIN_ERROR'
          });
          if (loginErrorChart) {
            keycloakMetrics.logins.error =
              loginErrorChart.data.columns[1][
                loginErrorChart.data.columns[1].length - 1
              ];
          }
          var registrationSuccessChart = _.findWhere(charts, {
            title: 'REGISTER'
          });
          if (registrationSuccessChart) {
            keycloakMetrics.registrations.success =
              registrationSuccessChart.data.columns[1][
                registrationSuccessChart.data.columns[1].length - 1
              ];
          }
          var registrationErrorChart = _.findWhere(charts, {
            title: 'REGISTER_ERROR'
          });
          if (registrationErrorChart) {
            keycloakMetrics.registrations.error =
              registrationErrorChart.data.columns[1][
                registrationErrorChart.data.columns[1].length - 1
              ];
          }

          keycloakMetrics.logins.total =
            keycloakMetrics.logins.success + keycloakMetrics.logins.error;
          keycloakMetrics.registrations.total =
            keycloakMetrics.registrations.success +
            keycloakMetrics.registrations.error;
          $scope.keycloakMetrics = keycloakMetrics;
        } else {
          $timeout(
            () => {
              // generic rendering of all data as line charts
              charts.forEach(chart => {
                c3.generate(chart);
              });
            },
            0,
            false
          );
        }
      }
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
  }
]);
