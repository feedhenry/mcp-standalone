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

    $scope.mobileapps = [];
    $scope.services = [];
    $scope.mcpError = false;

    ProjectsService.get($routeParams.project).then(
      _.spread(function(project, context) {
        $scope.project = project;
        $scope.projectContext = context;
        DataService.list(
          {
            group: 'servicecatalog.k8s.io',
            resource: 'serviceclasses'
          },
          context,
          function(serviceClasses) {
            $scope.serviceClasses = serviceClasses._data;
          }
        );
        mcpApi
          .mobileApps()
          .then(apps => {
            $scope.mobileapps = apps;
          })
          .catch(e => {
            console.error(e);
            $scope.mcpError = true;
          });

        mcpApi
          .mobileServices()
          .then(s => {
            $scope.services = s;
          })
          .catch(e => {
            console.error('error getting services ', e);
            $scope.mcpError = true;
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

    $scope.getIcon = function(service) {
      for (var serviceName in $scope.serviceClasses) {
        var serviceClass = $scope.serviceClasses[serviceName];
        if (
          serviceName === service.name ||
          serviceName.toLowerCase().indexOf(service.name) >= 0
        ) {
          if (
            typeof serviceClass.externalMetadata[
              'console.openshift.io/iconClass'
            ] !== 'undefined'
          ) {
            return formatIconClasses(
              serviceClass.externalMetadata['console.openshift.io/iconClass']
            );
          }
        }
      }
      return formatIconClasses('fa-clone');
    };

    formatIconClasses = function(icon) {
      bits = icon.split('-', 2);
      switch (bits[0]) {
        case 'font':
        case 'icon':
          return 'font-icon ' + icon;
        case 'fa':
          return 'fa ' + icon;
        default:
          return icon;
      }
    };
  }
]);
