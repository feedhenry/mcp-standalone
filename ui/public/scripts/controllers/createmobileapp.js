'use strict';

/**
 * @ngdoc function
 * @name mcp.controller:CreateMobileappController
 * @description
 * # CreateMobileappController
 * Controller of the mobileControlPanelApp
 */

angular
  .module('mobileControlPanelApp')
  .controller('CreateMobileappController', [
    '$scope',
    '$routeParams',
    '$location',
    'mcpApi',
    function($scope, $routeParams, $location, mcpApi) {
      $scope.alerts = {};
      $scope.projectName = $routeParams.project;

      $scope.breadcrumbs = [
        {
          title: $scope.projectName,
          link: 'project/' + $scope.projectName
        },
        {
          title: 'Mobile',
          link: 'project/' + $scope.projectName + '/mobile'
        },
        {
          title: 'Create Mobile App'
        }
      ];

      $scope.app = { clientType: '' };
      $scope.createApp = function() {
        mcpApi
          .createMobileApp($scope.app)
          .then(app => {
            $location.path(
              'project/' + $routeParams.project + '/browse/mobileapps'
            );
          })
          .catch(err => {
            console.error('failed to create app ', err);
          });
      };
    }
  ]);
