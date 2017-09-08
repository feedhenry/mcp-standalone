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
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        return $http.get(getMobileAppsURL(), requestConfig).then(res => {
          return res.data;
        });
      },
      mobileApp: function(id) {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        return $http
          .get(getMobileAppsURL() + '/' + id, requestConfig)
          .then(res => {
            return res.data;
          });
      },
      createMobileApp: function(mobileApp) {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        return $http
          .post(getMobileAppsURL(), mobileApp, requestConfig)
          .then(res => {
            return res.data;
          });
      },
      mobileServices: function() {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        return $http.get(getMobileServicesURL(), requestConfig).then(res => {
          return res.data;
        });
      },
      mobileService: function(name, withIntegrations) {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
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
      createMobileService: function(mobileService) {
        return $http
          .post(getMobileServicesURL(), mobileService, requestConfig)
          .then(res => {
            return res.data;
          });
      },
      integrateService: function(params) {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        let url = getMobileServicesURL() + '/configure';
        return $http.post(url, params, requestConfig).then(res => {
          return res.data;
        });
      }
    };
  }
]);
