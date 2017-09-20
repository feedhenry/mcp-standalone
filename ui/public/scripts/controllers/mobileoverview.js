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
  'AuthorizationService',
  function(
    $scope,
    $routeParams,
    $location,
    DataService,
    ProjectsService,
    mcpApi,
    AuthorizationService
  ) {
    Object.assign($scope, {
      projectName: $routeParams.project,
      alerts: {},
      renderOptions: Object.assign($scope.renderOptions || {}, {
        hideFilterWidget: true
      }),
      mcpError: false,
      overviews: {
        apps: {},
        services: {}
      }
    });

    ProjectsService.get($routeParams.project)
      .then(projectInfo => {
        const [project = {}, projectContext = {}] = projectInfo;
        $scope.project = project;
        $scope.projectContext = projectContext;

        $scope.overviews.apps = {
          title: 'Mobile Apps',
          actions: [
            {
              label: 'Create Mobile App',
              primary: true,
              action: $location.path.bind(
                $location,
                `project/${projectContext.projectName}/create-mobileapp`
              ),
              canView: function() {
                return true;
              }
            }
          ]
        };
        $scope.overviews.services = {
          title: 'Mobile Enabled Services',
          actions: [
            {
              label: 'Add External Service',
              action: $location.path.bind(
                $location,
                `project/${projectContext.projectName}/create-mobileservice`
              ),
              canView: function() {
                return AuthorizationService.canI(
                  'services',
                  'create',
                  projectContext.projectName
                );
              }
            },
            {
              label: 'Provision Catalog Service',
              primary: true,
              action: $location.path.bind($location, `/`),
              canView: function() {
                return AuthorizationService.canI(
                  'services',
                  'create',
                  projectContext.projectName
                );
              }
            }
          ]
        };

        return Promise.all([
          mcpApi.mobileApps(),
          mcpApi.mobileServices(),
          DataService.list(
            {
              group: 'servicecatalog.k8s.io',
              resource: 'serviceclasses'
            },
            $scope.projectContext
          )
        ]);
      })
      .then(overview => {
        const [apps = [], services = [], serviceClasses] = overview;
        $scope.overviews.apps.objects = apps;
        $scope.overviews.services.objects = services;
        $scope.overviews.services.serviceClasses = serviceClasses['_data'];
      })
      .catch(err => {
        console.error('Error getting overview ', err);
        $scope.mcpError = true;
      });

    $scope.actionSelected = function(object) {
      const ojectIsService = !!object.integrations;
      const actionFn = ojectIsService ? 'deleteService' : 'deleteApp';
      const getFn = ojectIsService ? 'mobileServices' : 'mobileApps';
      const objectType = ojectIsService ? 'services' : 'apps';
      mcpApi[actionFn](object)
        .then(result => {
          return mcpApi[getFn]();
        })
        .then(objects => {
          $scope.overviews[objectType].objects = objects;
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
