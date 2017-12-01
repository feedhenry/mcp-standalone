'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ServiceIntegrationsService
 * @description
 * # ServiceIntegrationsService
 * ServiceIntegrationsService
 */
angular.module('mobileControlPanelApp').service('ServiceIntegrationsService', [
  'McpService',
  function(McpService) {
    this.getIntegrationInfo = function(serviceId) {
      return Promise.all([
        McpService.mobileService(serviceId, 'true'),
        McpService.mobileApps()
      ]);
    };

    this.enableIntegration = function(integration) {
      return McpService.integrateService(integration);
    };
    this.disableIntegration = function(integration) {
      return McpService.deintegrateService(integration);
    };
  }
]);
