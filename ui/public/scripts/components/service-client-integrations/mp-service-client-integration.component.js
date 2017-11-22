'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-service-client-integration
 * @description
 * # mp-service-client-integration
 */
angular
  .module('mobileControlPanelApp')
  .component('mpServiceClientIntegration', {
    template: `<mp-client-integrations service=$ctrl.service apps=$ctrl.apps></mp-client-integrations>`,
    controller: [
      'ServiceClientService',
      '$routeParams',
      function(ServiceClientService, $routeParams) {
        ServiceClientService.getServiceClientInfo($routeParams.service)
          .then(serviceClientInfo => {
            const [service = {}, apps = []] = serviceClientInfo;

            this.service = service;
            this.apps = apps;
          })
          .catch(err => {
            console.error('Error loading Service client info', err);
          });
      }
    ]
  });
