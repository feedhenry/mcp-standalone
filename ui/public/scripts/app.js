'use strict';

console.log('MCP Extension Loaded');

// Add 'Mobile' to the left nav
window.OPENSHIFT_CONSTANTS.PROJECT_NAVIGATION.splice(1, 0, {
  label: 'Mobile',
  iconClass: 'fa fa-mobile',
  href: '/browse/mobileoverview',
  prefixes: [
    '/browse/mobileapps',
    '/browse/mobileservices',
    '/create-mobileapp',
    '/create-mobileservice'
  ],
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
      // this is reset to null so that it is not persisted from a request in one namespace
      // to a request in a different namespace
      window.MCP_URL = null;
      //GLOBAL_MCP_URL is used either for local development or a global MCP server
      if (typeof window.GLOBAL_MCP_URL !== 'undefined') {
        window.MCP_URL = window.GLOBAL_MCP_URL;
        return;
      }
      return ProjectsService.get($route.current.params.project).then(
        _.spread(function(project, context) {
          return DataService.get('routes', 'mcp-standalone', context, {
            errorNotification: false
          })
            .then(function(route) {
              window.MCP_URL =
                (route.spec.tls ? 'https://' : 'http://') + route.spec.host;
            })
            .catch(function(err) {
              // Ignore this error as the MCP_URL will be checked when needed,
              // and handle the absence of the url then
            });
        })
      );
    }
  ]
};

angular.module('mobileControlPanelApp', ['openshiftConsole']).config([
  '$routeProvider',
  function($routeProvider) {
    $routeProvider
      .when('/project/:project/create-mobileapp', {
        templateUrl: 'extensions/mcp/views/create-mobileapp.html',
        controller: 'CreateMobileappController',
        resolve: resolveMCPRoute
      })
      .when('/project/:project/create-mobileservice', {
        templateUrl: 'extensions/mcp/views/create-service.html',
        controller: 'CreateMobileServiceController',
        resolve: resolveMCPRoute
      })
      .when('/project/:project/browse/mobileoverview', {
        templateUrl: 'extensions/mcp/views/mobileoverview.html',
        controller: 'MobileOverviewController',
        reloadOnSearch: false,
        resolve: resolveMCPRoute
      })
      .when('/project/:project/browse/mobileapps/:mobileapp', {
        templateUrl: 'extensions/mcp/views/mobileapp.html',
        controller: 'MobileAppController',
        reloadOnSearch: false,
        resolve: resolveMCPRoute
      })
      .when('/project/:project/browse/mobileservices/:service', {
        templateUrl: 'extensions/mcp/views/mobileservice.html',
        controller: 'MobileServiceController',
        resolve: resolveMCPRoute
      });
  }
]);
