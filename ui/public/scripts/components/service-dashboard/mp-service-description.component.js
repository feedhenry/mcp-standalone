'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-service-description
 * @description
 * # mp-service-description
 */
angular.module('mobileControlPanelApp').component('mpServiceDescription', {
  template: `<div class="col-12">
                <h2 class="dashboard-details">Details</h2>
                <dl class="dl-horizontal left">
                  <dt>Description:</dt>
                  <dd>{{ $ctrl.getDescription($ctrl.service) }}</dd>
                  <dt>Dashboard URL:</dt>
                  <dd><a href="{{ $ctrl.service.host }}">{{$ctrl.service.host}}</a></dd>
                  <dt ng-repeat-start="(key, value) in $ctrl.service.params">{{ key }}</dt>
                  <dd ng-repeat-end>{{ value }}</dd>
                </dl>
              </div>`,
  bindings: {
    service: '<',
    serviceClasses: '<'
  },
  controller: [
    function() {
      this.getDescription = function(service) {
        if (!service) {
          return;
        }

        for (var scId in this.serviceClasses) {
          var sc = this.serviceClasses[scId];
          if (
            sc.spec.externalMetadata.hasOwnProperty('serviceName') &&
            sc.spec.externalMetadata.serviceName.toLowerCase() ==
              service.type.toLowerCase()
          ) {
            return sc.spec.description;
          }
        }
        return '';
      };
    }
  ]
});
