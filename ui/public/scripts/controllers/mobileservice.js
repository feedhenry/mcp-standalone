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
  'mcpApi',
  '$routeParams',
  function($scope, mcpApi, $routeParams) {
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

    $scope.integrations = [];
    mcpApi
      .mobileService($routeParams.service, 'true')
      .then(s => {
        $scope.service = s;
        $scope.integrations = Object.keys(s.integrations);
      })
      .catch(e => {
        console.error('failed to read service ', e);
      });
    mcpApi
      .mobileApps()
      .then(apps => {
        $scope.mobileappsCount = apps.length;
        $scope.mobileapps = {};
        for (var i = 0; i < apps.length; i++) {
          let app = apps[i];
          $scope.mobileapps[app.clientType] = 'true';
        }
        $scope.clients = Object.keys($scope.mobileapps);
        $scope.clientType = $scope.clients[0];
      })
      .catch(e => {
        console.error(e);
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
