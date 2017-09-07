'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.mcpApi
 * @description
 * # mcpApi
 * Service in the mobileControlPanelApp.
 */
angular.module('mobileControlPanelApp').service('mcpApi', [
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
      },
      mobileService: function(name, withIntegrations) {
        let url = getMobileServicesURL() + '/' + name;
        if (withIntegrations) {
          console.log('withIntegrations');
          url += '?withIntegrations=true';
        }
        console.log('calling ', url);
        return $http.get(url, requestConfig).then(res => {
          return res.data;
        });
      },
      integrateService: function(params) {
        let url = getMobileServicesURL() + '/configure';
        return $http.post(url, params, requestConfig).then(res => {
          return res.data;
        });
      }
    };
  }
]);
