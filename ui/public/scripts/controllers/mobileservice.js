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
    $scope.enabled = function(integration, service) {
      if (!service) {
        return false;
      }
      if (service.integrations[integration]) {
        return service.integrations[integration].enabled == true;
      }
      return false;
    };
    $scope.enableIntegration = function(service) {
      mcpApi
        .integrateService(service)
        .then(res => {
          console.log('Service integrated');
        })
        .catch(e => {
          console.log('error integrating service ', e);
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
