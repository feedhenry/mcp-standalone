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
  '$location',
  '$routeParams',
  '$filter',
  'ProjectsService',
  'mcpApi',
  'ServiceClassService',
  'DataService',
  'BuildsService',
  function(
    $scope,
    $location,
    $routeParams,
    $filter,
    ProjectsService,
    mcpApi,
    ServiceClassService,
    DataService,
    BuildsService
  ) {
    $scope.projectName = $routeParams.project;
    $scope.alerts = {};
    $scope.renderOptions = $scope.renderOptions || {};
    $scope.renderOptions.hideFilterWidget = true;
    $scope.breadcrumbs = [
      {
        title: 'Mobile App',
        link: 'project/' + $routeParams.project + '/browse/mobileoverview'
      },
      {
        title: $routeParams.mobileapp
      }
    ];

    $scope.installType = '';
    $scope.route = window.MCP_URL;

    const watches = [];
    const MOBILE_CI_CD_NAME = 'aerogear-digger';
    $scope.loading = true;
    $scope.dropdownActions = [
      {
        label: 'Edit',
        value: 'edit'
      }
    ];
    $scope.setView = function(view) {
      if (view === 'edit') {
        $location.url(
          `project/${$routeParams.project}/browse/mobileapps/${$routeParams.mobileapp}?tab=buildConfig`
        );
      }

      $scope.view = view;
    };
    $scope.startBuild = function() {
      BuildsService.startBuild($scope.buildConfig).then(() => {
        $location.url(
          `project/${$routeParams.project}/browse/mobileapps/${$routeParams.mobileapp}?tab=buildHistory`
        );
      });
    };

    $scope.createAppBuildConfig = function(appConfig) {
      appConfig.appID = $routeParams.mobileapp;
      mcpApi
        .createBuildConfig(appConfig)
        .then(response => {
          return DataService.get(
            'buildconfigs',
            appConfig.name,
            $scope.projectContext
          );
        })
        .then(res => {
          $scope.buildConfig = res;
          $scope.view = 'view';
        });
    };

    $scope.updateAppBuildConfig = function(appConfig) {
      DataService.update(
        'buildconfigs',
        appConfig.metadata.name,
        appConfig,
        $scope.projectContext
      )
        .then(() => {
          return DataService.get(
            'buildconfigs',
            appConfig.metadata.name,
            $scope.projectContext
          );
        })
        .then(res => {
          $scope.buildConfig = res;
          $scope.view = 'view';
        });
    };

    $scope.cancelEdit = function() {
      $scope.view = 'view';
    };

    var buildConfigForBuild = $filter('buildConfigForBuild');
    var filterBuilds = function(allBuilds) {
      $scope.builds = _.filter(allBuilds, build => {
        var buildConfigName = buildConfigForBuild(build) || '';
        return (
          $scope.buildConfig &&
          $scope.buildConfig.metadata.name === buildConfigName
        );
      });
      $scope.orderedBuilds = BuildsService.sortBuilds($scope.builds, true);
    };

    ProjectsService.get($routeParams.project)
      .then(function(projectInfo) {
        const [project = {}, projectContext = {}] = projectInfo;
        $scope.project = project;
        $scope.projectContext = projectContext;

        return Promise.all([
          DataService.list('buildconfigs', projectContext),
          DataService.list('builds', projectContext),
          DataService.list('secrets', projectContext),
          DataService.list(
            {
              group: 'servicecatalog.k8s.io',
              resource: 'clusterserviceclasses'
            },
            $scope.projectContext
          ),
          mcpApi.mobileApp($routeParams.mobileapp),
          mcpApi.mobileServices()
        ]);
      })
      .then(viewData => {
        const [
          buildConfigs = {},
          builds = {},
          secrets = {},
          serviceClasses,
          app = {},
          services = []
        ] = viewData;

        $scope.serviceClasses = serviceClasses['_data'];

        const buildData = buildConfigs['_data'];
        $scope.buildConfig = Object.keys(buildData)
          .map(key => {
            return buildData[key];
          })
          .filter(buildConfig => {
            return (
              buildConfig.metadata.labels['mobile-appid'] ===
              $routeParams.mobileapp
            );
          })
          .pop();

        $scope.view = $scope.buildConfig ? 'view' : 'create';

        filterBuilds(builds['_data']);

        watches.push(
          DataService.watch('builds', $scope.projectContext, function(builds) {
            filterBuilds(builds['_data']);
          })
        );

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

        $scope.integrations = services;
        $scope.hasMobileCiCd = Object.keys(secrets['_data'])
          .map(key => secrets['_data'][key])
          .some(secret => {
            return (
              secret.metadata.name === MOBILE_CI_CD_NAME &&
              secret.metadata.namespace === $scope.project.metadata.name
            );
          });

        $scope.loading = false;
      });

    $scope.installationOpt = function(type) {
      $scope.installType = type;
    };
    $scope.sample = 'code';
    $scope.codeOpts = function(type) {
      $scope.sample = type;
    };

    $scope.openServiceIntegration = function(id) {
      $location.url(
        `project/${$routeParams.project}/browse/mobileservices/${id}?tab=integrations`
      );
    };

    $scope.$on('$destroy', function() {
      DataService.unwatchAll(watches);
    });

    $scope.getServiceName = function(service) {
      const serviceClass = ServiceClassService.retrieveServiceClass(
        service,
        $scope.serviceClasses
      );
      return ServiceClassService.retrieveDisplayName(
        serviceClass,
        service.name
      );
    };

    $scope.getServiceIcon = function(service) {
      const serviceClass = ServiceClassService.retrieveServiceClass(
        service,
        $scope.serviceClasses
      );
      return ServiceClassService.retrieveIcon(serviceClass);
    };
  }
]);
