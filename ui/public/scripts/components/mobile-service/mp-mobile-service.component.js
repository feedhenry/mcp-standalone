'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-mobile-service
 * @description
 * # mp-mobile-service
 */
angular.module('mobileControlPanelApp').component('mpMobileService', {
  template: `<div class="mp-service">
              <div class="container-fluid">
                <mp-service-info service=$ctrl.service></mp-service-info>
              </div>
            </div>`,
  controller: [
    'McpService',
    '$routeParams',
    function(McpService, $routeParams) {
      this.$onInit = function() {
        McpService.mobileService($routeParams.service, 'true')
          .then(service => (this.service = service))
          .catch(err => console.error('Error loading Services', err));
      };
    }
  ]
});
