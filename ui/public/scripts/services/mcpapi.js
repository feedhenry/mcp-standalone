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
      deleteService: function(service) {
        if (!window.MCP_URL) {
          return Promise.reject('No MCP URL');
        }

        return $http
          .delete(`${getMobileServicesURL()}/${service.id}`, requestConfig)
          .then(res => {
            return res.data;
          });
      },
      deleteApp: function(app) {
        if (!window.MCP_URL) {
          return Promise.reject('No MCP URL');
        }

        return $http
          .delete(`${getMobileAppsURL()}/${app.id}`, requestConfig)
          .then(res => {
            return res.data;
          });
      },
      mobileServiceMetrics: function(name) {
        // TODO: Avoid duplicate code for checking existance of MCP Url
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        let url = getMobileServicesURL() + '/' + name + '/metrics';
        console.log('calling ', url);
        return $http.get(url, requestConfig).then(res => {
          // return new Promise(function(resolve, reject) {
          // var res = {
          //   data: [
          //     {
          //       x: [
          //         '2013-01-01 11:22:45',
          //         '2013-01-02 11:22:45',
          //         '2013-01-03 11:22:45',
          //         '2013-01-04 11:22:45',
          //         '2013-01-05 11:22:45',
          //         '2013-01-06 11:22:45'
          //       ],
          //       y: {
          //         data5: [90, 150, 160, 165, 180, 5]
          //       }
          //     },
          //     {
          //       x: [
          //         '2013-01-01 11:22:45',
          //         '2013-01-02 11:22:45',
          //         '2013-01-03 11:22:45',
          //         '2013-01-04 11:22:45',
          //         '2013-01-05 11:22:45',
          //         '2013-01-06 11:22:45'
          //       ],
          //       y: {
          //         data3: [70, 100, 390, 295, 170, 220]
          //       }
          //     }
          //   ]
          // };
          // console.log('res.data', res.data);

          var data = [];
          res.data.forEach(raw => {
            // x axis
            //   from -> ['val1', 'val2', val3'],
            //     to -> ['x', 'val1', 'val2', val3'],
            raw.x.unshift('x');

            _.forIn(raw.y, (val, key) => {
              var columns = [raw.x];

              // y axis
              //   from -> [30, 200, 100, 400, 150, 250]
              //     to -> ['data1', 30, 200, 100, 400, 150, 250],
              val.unshift(key);
              columns.push(val);
              data.push(columns);
            });
          });

          return data;
        });
      },
      integrateService: function(params) {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        let url =
          getMobileServicesURL() +
          '/configure/' +
          params.component +
          '/' +
          params.service;
        return $http.post(url, {}, requestConfig).then(res => {
          return res.data;
        });
      },
      deintegrateService: function(params) {
        if (!window.MCP_URL) {
          return new Promise(function(resolve, reject) {
            return reject('No MCP URL');
          });
        }
        let url =
          getMobileServicesURL() +
          '/configure/' +
          params.component +
          '/' +
          params.service;
        return $http.delete(url, requestConfig).then(res => {
          return res.data;
        });
      }
    };
  }
]);
