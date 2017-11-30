'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileServiceController
 * @description
 * # MobileServiceController
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp').controller('MobileServiceController', [
  '$scope',
  '$timeout',
  'McpService',
  '$routeParams',
  'ProjectsService',
  'DataService',
  function(
    $scope,
    $timeout,
    McpService,
    $routeParams,
    ProjectsService,
    DataService
  ) {
    $scope.alerts = {};
    $scope.projectName = $routeParams.project;
    $scope.breadcrumbs = [
      {
        title: 'Mobile Service',
        link: 'project/' + $routeParams.project + '/browse/mobileoverview'
      },
      {
        title: $routeParams.service
      }
    ];

    const knownServices = ['fh-sync-server', 'keycloak'];
    $scope.chartData = [];
    $scope.integrations = [];
    McpService.mobileServiceMetrics($routeParams.service)
      .then(chartData => {
        $scope.chartData = chartData;
        $timeout(() => {});
      })
      .catch(err => {
        console.error('Error loading Service Metrics', err);
      });
    ProjectsService.get($routeParams.project).then(projectInfo => {
      const [project = {}, projectContext = {}] = projectInfo;
      $scope.project = project;
      $scope.projectContext = projectContext;
      DataService.list(
        {
          group: 'servicecatalog.k8s.io',
          resource: 'clusterserviceclasses'
        },
        $scope.projectContext,
        function(serviceClasses) {
          $scope.serviceClasses = serviceClasses['_data'];
        }
      );
    });
    Promise.all([
      McpService.mobileService($routeParams.service, 'true'),
      McpService.mobileApps()
    ])
      .then(serviceInfo => {
        const [service = {}, apps = []] = serviceInfo;

        $scope.service = service;
        $scope.integrations = Object.keys(service.integrations);
        $scope.service.integrations = Object.keys(
          service.integrations
        ).map(key => {
          return Object.assign(service.integrations[key], { target: service });
        });

        $scope.templateName = knownServices.includes(service.name)
          ? service.name
          : 'default-service';

        $scope.mobileappsCount = apps.length;
        $scope.mobileapps = apps.reduce((acc, current) => {
          acc[current.clientType] = 'true';
          return acc;
        }, {});
        $scope.clients = Object.keys($scope.mobileapps);
        $scope.clientType = $scope.clients[0];
        $timeout(() => {});
      })
      .catch(err => {
        console.error('Error loading Services', err);
      });

    $scope.status = function(integration, service) {
      if ($scope.processing(integration, service)) {
        return 2;
      }
      if ($scope.enabled(integration, service)) {
        return 1;
      }
      return 0;
    };
    $scope.serviceWritable = function(service) {
      return service.writable;
    };
    $scope.enabled = function(integration, service) {
      if (!service) {
        return false;
      }
      if (service.integrations[integration]) {
        return service.integrations[integration].enabled === true;
      }
      return false;
    };
    $scope.processing = function(integration, service) {
      if (!service) {
        return false;
      }
      if (
        service.integrations[integration] &&
        service.integrations[integration].processing
      ) {
        return service.integrations[integration].processing === true;
      }
      return false;
    };

    $scope.integrationToggled = function(integration, enabled) {
      if (enabled) {
        $scope.enableIntegration(integration);
      } else {
        $scope.disableIntegration(integration);
      }
    };

    $scope.enableIntegration = function(integration) {
      McpService.integrateService(integration)
        .then(res => {
          //inspect res
          integration.enabled = true;
        })
        .catch(e => {
          console.log('error integrating service ', e);
        });
      return true;
    };
    $scope.disableIntegration = function(integration) {
      McpService.deintegrateService(integration)
        .then(res => {
          //inspect res
          integration.enabled = false;
        })
        .catch(e => {
          console.log('error deintegrating service ', e);
        });
      return true;
    };

    $scope.clientType =
      $scope.clients && $scope.clients.length > 0
        ? $scope.clients[0]
        : 'cordova';

    $scope.installationOpt = function(type) {
      $scope.clientType = type;
    };

    $scope.dashboardSelected = function() {
      window.dispatchEvent(new Event('resize'));
    };

    $scope.getDescription = function(service) {
      for (var scId in $scope.serviceClasses) {
        var sc = $scope.serviceClasses[scId];
        if (
          sc.spec.externalMetadata.hasOwnProperty('serviceName') &&
          sc.spec.externalMetadata.serviceName.toLowerCase() ==
            service.type.toLowerCase()
        ) {
          return sc.spec.description;
        }
      }
      return '';
    };
  }
]);
