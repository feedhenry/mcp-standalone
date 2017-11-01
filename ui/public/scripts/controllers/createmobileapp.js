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
          title: 'Mobile Apps',
          link: 'project/' + $routeParams.project + '/browse/mobileoverview'
        },
        {
          title: 'Create Mobile App'
        }
      ];

      $scope.created = function(app) {
        mcpApi
          .createMobileApp(app)
          .then(app => {
            $location.path(
              'project/' + $routeParams.project + '/browse/mobileoverview'
            );
          })
          .catch(err => {
            console.error('failed to create app ', err);
          });
      };

      $scope.cancelled = function() {
        $location.path(
          'project/' + $routeParams.project + '/browse/mobileoverview'
        );
      };
    }
  ]);
