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
      getStartedOptions: {
        modalOpen: false
      },
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
          type: 'app',
          title: 'Mobile Apps',
          text:
            'You can create a Mobile App to enable Mobile Integrations with Mobile Enabled Services',
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
          type: 'service',
          title: 'Mobile Enabled Services',
          modalOpen: false,
          text:
            'You can provision or link a Mobile Enabled Service to enable a Mobile App Integration.',
          actions: [
            {
              label: 'Add External Service',
              modal: true,
              contentUrl:
                'extensions/mcp/templates/create-service.template.html',
              action: function(err) {
                if (err) {
                  return;
                }

                mcpApi.mobileServices().then(services => {
                  $scope.overviews.services.objects = services;
                  $scope.overviews.services.modalOpen = false;
                });
              },
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
              resource: 'clusterserviceclasses'
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
      const objectIsService = !!object.integrations;
      const actionFn = objectIsService ? 'deleteService' : 'deleteApp';
      const getFn = objectIsService ? 'mobileServices' : 'mobileApps';
      const objectType = objectIsService ? 'services' : 'apps';
      mcpApi[actionFn](object)
        .then(result => {
          return mcpApi[getFn]();
        })
        .then(objects => {
          $scope.overviews[objectType].objects = objects;
        });
    };

    $scope.serviceCreated = function(err) {
      if (err) {
        return;
      }

      return mcpApi.mobileServices().then(services => {
        $scope.overviews.services.objects = services;
        $scope.getStartedOptions.modalOpen = false;
      });
    };

    $scope.objectSelected = function(object) {
      const objectIsService = !!object.integrations;
      if (objectIsService) {
        $location.path(
          'project/' +
            $routeParams.project +
            '/browse/mobileservices/' +
            object.id
        );
      } else {
        $location.path(
          'project/' + $routeParams.project + '/browse/mobileapps/' + object.id
        );
      }
    };
  }
]);
