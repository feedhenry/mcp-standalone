'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-chart-data
 * @description
 * # mp-chart-data
 */
angular.module('mobileControlPanelApp').component('mpChartData', {
  template: `<div class="col-12">
                <div ng-if="$ctrl.templateName" ng-include="'extensions/mcp/templates/' + $ctrl.templateName + '-chart.template.html'"></div>
              </div>

              <div ng-if="$ctrl.chartData.length === 0">
                <div class="empty-state-message text-center">
                  <h2>No Stats currently available for this service</h2>
                </div>
              </div>`,
  bindings: {
    chartData: '<',
    service: '<'
  },
  controller: [
    function() {
      const knownServices = ['fh-sync-server', 'keycloak'];
      this.$onChanges = function(changes) {
        const service = changes.service && changes.service.currentValue;
        if (!service) {
          return;
        }

        this.templateName = knownServices.includes(service.name)
          ? service.name
          : 'default-service';
      };
    }
  ]
});
