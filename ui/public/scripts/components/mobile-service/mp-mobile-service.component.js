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
                <mp-service-info service=$ctrl.service service-classes=$ctrl.serviceClasses></mp-service-info>
              </div>
            </div>`,
  controller: [
    'MobileServiceService',
    '$routeParams',
    function(MobileServiceService, $routeParams) {
      this.$onInit = function() {
        MobileServiceService.getData($routeParams.project, $routeParams.service)
          .then(data => {
            const [service = {}, serviceClasses = []] = data;
            this.service = service;
            this.serviceClasses = serviceClasses['_data'];
          })
          .catch(err => console.error('Error loading Services', err));
      };
    }
  ]
});
