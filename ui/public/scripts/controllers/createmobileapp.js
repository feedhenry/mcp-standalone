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
    'DataService',
    'ProjectsService',
    'APIService',
    function(
      $scope,
      $routeParams,
      $location,
      DataService,
      ProjectsService,
      APIService
    ) {
      $scope.alerts = {};
      $scope.projectName = $routeParams.project;

      $scope.breadcrumbs = [
        {
          title: 'Overview',
          link: 'project/' + $routeParams.project + '/overview'
        },
        {
          title: 'Create Mobile App'
        }
      ];

      /*
      configMap = {
        apiVersion: 'v1',
        kind: 'ConfigMap',
        metadata: {
          namespace: $routeParams.project
        },
        data: {}
      };
      */
      var uuidv4 = function() {
        var uuid = '',
          i,
          random;
        for (i = 0; i < 32; i++) {
          random = (Math.random() * 16) | 0;

          if (i == 8 || i == 12 || i == 16 || i == 20) {
            uuid += '-';
          }
          uuid += (i == 12 ? 4 : i == 16 ? (random & 3) | 8 : random).toString(
            16
          );
        }
        return uuid;
      };

      var convertMobileAppToConfigMap = function(app) {
        return {
          kind: 'ConfigMap',
          apiVersion: 'v1',
          metadata: {
            name: app.name + '-' + Date.now(),
            namespace: $routeParams.project,
            labels: {
              group: 'mobileapp',
              name: app.name
            }
          },
          data: {
            name: app.name,
            displayName: app.displayName,
            clientType: app.clientType,
            apiKey: uuidv4(),
            description: app.description
          }
        };
      };

      /*
      {
        "clientType": "android",
        "name": "aa",
        "displayName": "aa",
        "description": "aa"
      }
      */
      $scope.created = function(app) {
        ProjectsService.get($routeParams.project).then(
          _.spread(function(project, context) {
            $scope.project = project;

            var configMap = convertMobileAppToConfigMap(app);
            var createConfigMapVersion = APIService.objectToResourceGroupVersion(
              configMap
            );
            DataService.create(
              createConfigMapVersion,
              null,
              configMap,
              context
            ).then(
              function() {
                // Success
                $location.path('project/' + $routeParams.project + '/overview');
              },
              function(err) {
                // Failure
                console.error('failed to create app ', err);
              }
            );
          })
        );
      };

      $scope.cancelled = function() {
        $location.path('project/' + $routeParams.project + '/overview');
      };
    }
  ]);
