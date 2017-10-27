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
  }
]);
