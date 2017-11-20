'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-client-setup
 * @description
 * # mp-client-setup
 */
angular.module('mobileControlPanelApp').component('mpClientSetup', {
  template: `<div>
              <h3>Core SDK Setup</h3>
              <p>
                {{$ctrl.app.Description}}
              </p>
              <p>
                The core sdk will discover the services you have provisioned and make them available to your application by exposing the
                required configuration. This will allow you to easily and simply plug in and start using your chosen mobile services.
              </p>

              <mp-android-setup app=$ctrl.app></mp-android-setup>
              <mp-cordova-setup app=$ctrl.app></mp-cordova-setup>
              <mp-ios-setup app=$ctrl.app></mp-ios-setup>

              <mp-service-integrations
                integrations=$ctrl.integrations
                service-classes=$ctrl.serviceClasses
                service-selected=$ctrl.openServiceIntegration>
              </mp-service-integrations>
            </div>`,
  controller: [
    '$routeParams',
    '$location',
    'ClientSetupService',
    function($routeParams, $location, ClientSetupService) {
      this.$onInit = function() {
        ClientSetupService.getData(
          $routeParams.project,
          $routeParams.mobileapp
        ).then(data => {
          const [app = {}, services = [], serviceClasses = []] = data;
          this.app = app;
          this.integrations = services;
          this.serviceClasses = serviceClasses['_data'];
        });
      };

      this.openServiceIntegration = function(id) {
        $location.url(
          `project/${$routeParams.project}/browse/mobileservices/${id}?tab=integrations`
        );
      };
    }
  ]
});
