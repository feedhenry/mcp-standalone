'use strict';

console.log('MCP Extension Loaded');

// Add 'Mobile' to the left nav
window.OPENSHIFT_CONSTANTS.PROJECT_NAVIGATION.splice(1, 0, {
  label: 'Mobile',
  iconClass: 'fa fa-mobile',
  href: '/browse/mobileoverview',
  prefixes: ['/browse/mobileoverview'],
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

angular.module('mobileControlPanelApp', ['openshiftConsole']).config([
  '$routeProvider',
  function($routeProvider) {
    $routeProvider
      .when('/project/:project/create-mobileapp', {
        templateUrl: 'extensions/mcp/views/create-mobileapp.html',
        controller: 'CreateMobileappController',
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
