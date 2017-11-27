'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ClientSetupService
 * @description
 * # ClientSetupService
 * ClientSetupService
 */
angular.module('mobileControlPanelApp').service('ClientSetupService', [
  'McpService',
  function(McpService) {
    this.getData = function(appId) {
      return Promise.all([
        McpService.mobileApp(appId),
        McpService.mobileServices()
      ]);
    };
  }
]);
