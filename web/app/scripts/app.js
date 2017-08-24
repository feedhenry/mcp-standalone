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
  .config(['$locationProvider', '$routeProvider', 'RedirectLoginServiceProvider', function ($locationProvider, $routeProvider, RedirectLoginServiceProvider) {
    $locationProvider.html5Mode(true);

    $routeProvider
      .when('/mobileapps', {
        templateUrl: 'views/mobileapps.html',
        controller: 'MobileappsCtrl',
        requireAuthentication: true
      })
      .when('/mobileapp', {
        templateUrl: 'views/mobileapp.html',
        controller: 'MobileappCtrl',
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
      .otherwise({
        redirectTo: '/mobileapps'
      });

      RedirectLoginServiceProvider.OAuthScope('user:info user:check-access');
  }])
  .filter('debug', function() {
    return function(input) {
      if (input === '') return 'empty string';
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

