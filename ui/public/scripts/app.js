'use strict';

console.log('MCP Extension Loaded');

// Add 'Mobile' to the left nav
window.OPENSHIFT_CONSTANTS.PROJECT_NAVIGATION.splice(1, 0, {
  label: 'Mobile',
  iconClass: 'fa fa-mobile',
  href: '/browse/mobileapps',
  prefixes: ['/browse/mobileapps'],
  isValid: function() {
    // TODO: Can this check if any mobile apps exist first?
    return true;
  }
});

// Add 'Mobile' category and sub-categories to the Service Catalog UI
window.OPENSHIFT_CONSTANTS.SERVICE_CATALOG_CATEGORIES.splice(
  OPENSHIFT_CONSTANTS.SERVICE_CATALOG_CATEGORIES.length,
  0,
  {
    id: 'mobile',
    label: 'Mobile',
    subCategories: [
      { id: 'apps', label: 'Apps', tags: ['mobile'], icon: 'fa fa-mobile' },
      {
        id: 'services',
        label: 'Services',
        tags: ['mobile-service'],
        icon: 'fa fa-database'
      }
    ]
  }
);

var resolveMCPRoute = {
  MCPRoute: [
    '$route',
    'ProjectsService',
    'DataService',
    function($route, ProjectsService, DataService) {
      if (window.MCP_URL) {
        return;
      }
      return ProjectsService.get($route.current.params.project).then(
        _.spread(function(project, context) {
          return DataService.get('routes', 'mcp-standalone', context, {
            errorNotification: false
          }).then(function(route) {
            window.MCP_URL =
              (route.spec.tls ? 'https://' : 'http://') + route.spec.host;
          });
        })
      );
    }
  ]
};

angular
  .module('mobileControlPanelApp', ['openshiftConsole'])
  .config([
    '$routeProvider',
    function($routeProvider) {
      $routeProvider
        .when('/project/:project/create-mobileapp', {
          templateUrl: 'extensions/mcp/views/create-mobileapp.html',
          controller: 'CreateMobileappController',
          resolve: resolveMCPRoute
        })
        .when('/project/:project/browse/mobileapps', {
          templateUrl: 'extensions/mcp/views/mobileapps.html',
          controller: 'MobileAppsController',
          reloadOnSearch: false,
          resolve: resolveMCPRoute
        })
        .when('/project/:project/browse/mobileapps/:mobileapp', {
          templateUrl: 'extensions/mcp/views/mobileapp.html',
          controller: 'MobileAppController',
          reloadOnSearch: false,
          resolve: resolveMCPRoute
        });
    }
  ])
  .service('mcpApi', [
    '$http',
    'AuthService',
    function($http, AuthService) {
      function getMobileAppsURL() {
        return window.MCP_URL + '/mobileapp';
      }
      function getMobileServicesURL() {
        return window.MCP_URL + '/mobileservice';
      }
      // AngularJS will instantiate a singleton by calling "new" on this function
      let requestConfig = { headers: {} };
      AuthService.addAuthToRequest(requestConfig);

      return {
        mobileApps: function() {
          return $http.get(getMobileAppsURL(), requestConfig).then(res => {
            return res.data;
          });
        },
        mobileApp: function(id) {
          return $http
            .get(getMobileAppsURL() + '/' + id, requestConfig)
            .then(res => {
              return res.data;
            });
        },
        createMobileApp: function(mobileApp) {
          return $http
            .post(getMobileAppsURL(), mobileApp, requestConfig)
            .then(res => {
              return res.data;
            });
        },
        mobileServices: function() {
          return $http.get(getMobileServicesURL(), requestConfig).then(res => {
            return res.data;
          });
        }
      };
    }
  ])
  .controller('MobileAppsController', [
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
    }
  ])
  .controller('MobileAppController', [
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
  ])
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
          title: $scope.projectName,
          link: 'project/' + $scope.projectName
        },
        {
          title: 'Mobile',
          link: 'project/' + $scope.projectName + '/mobile'
        },
        {
          title: 'Create Mobile App'
        }
      ];

      $scope.app = { clientType: '' };
      $scope.createApp = function() {
        mcpApi
          .createMobileApp($scope.app)
          .then(app => {
            $location.path(
              'project/' + $routeParams.project + '/browse/mobileapps'
            );
          })
          .catch(err => {
            console.error('failed to create app ', err);
          });
      };
    }
  ]);

hawtioPluginLoader.addModule('mobileControlPanelApp');
