'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileOverviewController
 * @description
 * # MobileOverviewController
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp').controller('MobileOverviewController', [
  '$scope',
  '$routeParams',
  '$location',
  'DataService',
  'ProjectsService',
  'mcpApi',
  function(
    $scope,
    $routeParams,
    $location,
    DataService,
    ProjectsService,
    mcpApi
  ) {
    $scope.projectName = $routeParams.project;
    $scope.alerts = {};
    $scope.renderOptions = $scope.renderOptions || {};
    $scope.renderOptions.hideFilterWidget = true;
    $scope.breadcrumbs = [
      {
        title: 'Mobile',
        link: 'project/' + $routeParams.project + '/browse/mobileapps'
      },
      {
        title: $routeParams.mobileapp
      }
    ];

    $scope.mobileapps = [];
    $scope.services = [];

    ProjectsService.get($routeParams.project).then(
      _.spread(function(project, context) {
        $scope.project = project;
        $scope.projectContext = context;
        mcpApi
          .mobileApps()
          .then(apps => {
            $scope.mobileapps = apps;
          })
          .catch(e => {
            console.error(e);
          });

        mcpApi
          .mobileServices()
          .then(s => {
            $scope.services = s;
          })
          .catch(e => {
            console.error('error getting services ', e);
          });
      })
    );

    $scope.openApp = function(id) {
      $location.path(
        'project/' + $routeParams.project + '/browse/mobileapps/' + id
      );
    };

    $scope.openService = function(id) {
      $location.path(
        'project/' + $routeParams.project + '/browse/mobileservices/' + id
      );
    };
  }
]);
