'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ServiceClientService
 * @description
 * # ServiceClientService
 * ServiceClientService
 */
angular.module('mobileControlPanelApp').service('ServiceClientService', [
  'McpService',
  function(McpService) {
    this.getServiceClientInfo = function(serviceId) {
      return Promise.all([
        McpService.mobileService(serviceId, 'true'),
        McpService.mobileApps()
      ]);
    };

    this.enabled = function(integration, service) {
      if (!service) {
        return false;
      }
      if (service.integrations[integration]) {
        return service.integrations[integration].enabled === true;
      }
      return false;
    };
  }
]);
