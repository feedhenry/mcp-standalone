'use strict';

/**
 * @ngdoc overview
 * @name mobileControlPanelApp
 * @description
 * # mobileControlPanelApp
 *
 * Main module of the application.
 */
angular
  .module('mobileControlPanelApp', [
    'ngRoute',
    'openshiftCommonServices'
  ])
  .config(['$routeProvider', '$locationProvider','RedirectLoginServiceProvider', function ($routeProvider, $locationProvider ,RedirectLoginServiceProvider) {
    $locationProvider.html5Mode(true);
    $routeProvider
      .when('/apps', {
        templateUrl: '/views/mobileapps.html',
        //todo change the name of the controller
        controller: 'MobileappsCtrl',
        requireAuthentication: true
      })
      .when('/apps/:id', {
        templateUrl: 'views/mobileapp.html',
        controller: 'MobileappCtrl',
        requireAuthentication: true
      })
      .when('/appcreate', {
        templateUrl: '/views/appcreate.html',
        controller: 'AppCreateCtrl',
        requireAuthentication: true
      })
      .when('/oauth', {
        templateUrl: 'views/oauth.html',
        controller: 'OauthCtrl'
      })
      .when('/error', {
        templateUrl: 'views/error.html',
        controller: 'ErrorCtrl'
      })
      .when('/integrations', {
        templateUrl: 'views/integrations.html',
        controller: 'IntegrationsCtrl',
        controllerAs: 'integrations'
      })
      .when('/integrations/:service', {
        templateUrl: 'views/integration.html',
        controller: 'IntegrationCtrl',
        controllerAs: 'integration'
      })
      .otherwise({
        redirectTo: '/apps'
      });

      RedirectLoginServiceProvider.OAuthScope(window.OPENSHIFT_CONFIG.auth.scope);
  }])
  .filter('debug', function() {
    return function(input) {
      if (input === ''){ return 'empty string';}
      return input ? input : ('' + input);
    };
  })
  .run(['$rootScope', '$location', 'AuthService', function ($rootScope, $location, AuthService) {
    $rootScope.$on('$routeChangeStart', function (event, url) {
      if (url.requireAuthentication) {
        AuthService.withUser().then(function() {
          // no further action. Login check was successful
        });
      }
    });
}]);

hawtioPluginLoader.addModule('mobileControlPanelApp');