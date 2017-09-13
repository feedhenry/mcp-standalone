'use strict';

/**
 * @ngdoc function
 * @name mcp.controller:CreateServiceController
 * @description
 * # CreateServiceController
 * Controller of the mobileControlPanelApp
 */

angular
  .module('mobileControlPanelApp')
  .controller('CreateMobileServiceController', [
    '$scope',
    '$routeParams',
    'mcpApi',
    '$location',
    'DataService',
    function($scope, $routeParams, mcpApi, $location, DataService) {
      $scope.breadcrumbs = [
        {
          title: 'Mobile Apps',
          link: 'project/' + $routeParams.project + '/browse/mobileoverview'
        },
        {
          title: 'Create Mobile Service'
        }
      ];
      DataService.list('projects', {})
        .then(p => {
          $scope.projects = p._data;
        })
        .catch(e => {
          console.log('error listing projects', e);
        });
      $scope.projectName = $routeParams.project;
      $scope.customFields = [];
      $scope.externalService = {
        labels: { external: 'true' },
        params: {}
      };
      $scope.addService = function() {
        for (i = 0; i < $scope.customFields.length; i++) {
          var f = $scope.customFields[i];
          if (f.name != '' && f.value != '') {
            $scope.externalService.params[f.name] = f.value;
          }
        }
        console.log('called add service', $scope.externalService);
        mcpApi
          .createMobileService($scope.externalService)
          .then(s => {
            $location.path(
              'project/' + $routeParams.project + '/browse/mobileoverview'
            );
          })
          .catch(err => {
            console.error('failed to create mobile service ', err);
          });
      };
    }
  ]);
