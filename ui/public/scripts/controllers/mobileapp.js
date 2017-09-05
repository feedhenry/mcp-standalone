'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileAppController
 * @description
 * # MobileAppController
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp').controller('MobileAppController', [
  '$scope',
  '$routeParams',
  'ProjectsService',
  'mcpApi',
  function($scope, $routeParams, ProjectsService, mcpApi) {
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

    $scope.installType = '';
    $scope.route = window.MCP_URL;

    ProjectsService.get($routeParams.project).then(
      _.spread(function(project, context) {
        $scope.project = project;
        $scope.projectContext = context;
        mcpApi
          .mobileApp($routeParams.mobileapp)
          .then(app => {
            $scope.app = app;
            switch (app.clientType) {
              case 'cordova':
                $scope.installType = 'npm';
                break;
              case 'android':
                $scope.installType = 'maven';
                break;
              case 'iOS':
                $scope.installType = 'cocoapods';
                break;
            }
          })
          .catch(e => {
            console.error('failed to read app', e);
          });
      })
    );

    $scope.installationOpt = function(type) {
      $scope.installType = type;
    };
    $scope.sample = 'code';
    $scope.codeOpts = function(type) {
      $scope.sample = type;
    };
  }
]);
