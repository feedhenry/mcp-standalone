'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-service-info
 * @description
 * # mp-service-info
 */
angular.module('mobileControlPanelApp').component('mpServiceInfo', {
  template: `<div class="row">
                <div class="col-md-12">
                  <div class="mp-object-tab-title">
                    <h1>
                      {{$ctrl.service.name}}
                    </h1>
                  </div>
                  <uib-tabset class="mp-tabsets" justified=true>
                    <uib-tab ng-if="$ctrl.service" active=true select=$ctrl.dashboardSelected()>
                      <uib-tab-heading>Dashboard</uib-tab-heading>
                      <mp-service-dashboard></mp-service-dashboard>
                    </uib-tab>
                    <uib-tab>
                      <uib-tab-heading>App Integrations</uib-tab-heading>
                      <mp-service-client-integration></mp-service-client-integration>
                    </uib-tab>
                    <uib-tab ng-if="$ctrl.hasIntegrations">
                      <uib-tab-heading>Services Integrations</uib-tab-heading>
                      <mp-mobile-service-integration></mp-mobile-service-integration>
                    </uib-tab>
                  </uib-tabset>
                </div>
              </div>`,
  bindings: {
    service: '<'
  },
  controller: [
    function() {
      this.$onChanges = function(changes) {
        const service = changes.service && changes.service.currentValue;
        if (!service) {
          return;
        }

        this.hasIntegrations = !!Object.keys(service.integrations).length;
      };

      this.dashboardSelected = function() {
        window.dispatchEvent(new Event('resize'));
      };
    }
  ]
});
