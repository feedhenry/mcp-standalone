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
    Object.assign($scope, {
      projectName: $routeParams.project,
      alerts: {},
      renderOptions: Object.assign($scope.renderOptions || {}, {
        hideFilterWidget: true
      }),
      mcpError: false,
      serviceOptions: {},
      appOptions: {}
    });

    ProjectsService.get($routeParams.project)
      .then(projectInfo => {
        const [project = {}, projectContext = {}] = projectInfo;
        $scope.project = project;
        $scope.projectContext = projectContext;

        DataService.list(
          {
            group: 'servicecatalog.k8s.io',
            resource: 'serviceclasses'
          },
          $scope.projectContext,
          function(serviceClasses) {
            $scope.serviceClasses = serviceClasses._data;
          }
        );

        return Promise.all([mcpApi.mobileApps(), mcpApi.mobileServices()]);
      })
      .then(mobileOverview => {
        const [mobileapps = [], services = []] = mobileOverview;
        $scope.mobileapps = mobileapps;
        $scope.services = services;
      })
      .catch(err => {
        console.error('Error getting overview ', err);
        $scope.mcpError = true;
      });

    $scope.actionSelected = function(object) {
      const ojectIsService = !!object.integrations;
      const actionFn = ojectIsService ? 'deleteService' : 'deleteApp';
      const getFn = ojectIsService ? 'mobileServices' : 'mobileApps';
      const objectType = ojectIsService ? 'services' : 'mobileapps';
      mcpApi[actionFn](object)
        .then(result => {
          return mcpApi[getFn]();
        })
        .then(objects => {
          $scope[objectType] = objects;
        });
    };

    $scope.openApp = function(app) {
      $location.path(
        'project/' + $routeParams.project + '/browse/mobileapps/' + app.id
      );
    };

    $scope.openService = function(service) {
      $location.path(
        'project/' +
          $routeParams.project +
          '/browse/mobileservices/' +
          service.id
      );
    };
  }
]);
