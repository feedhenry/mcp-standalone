'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-service-integrations
 * @description
 * # mp-service-integrations
 */
angular.module('mobileControlPanelApp').component('mpServiceIntegrations', {
  template: `<div>
              <h3>Available Integrations</h3>

              <div ng-if="$ctrl.integrations.length > 0" class="row row-cards-pf">
                <div ng-repeat="service in $ctrl.integrations" class="col-xs-12 col-sm-4 col-md-4 col-lg-2">
                  <div class="card-pf card-pf-view card-pf-view-select card-pf-view-multi-select" ng-click="$ctrl.serviceSelected(service.id)">
                    <div class="card-pf-body">
                      <div class="card-pf-top-element">
                        <mp-object-icon object=service service-classes=$ctrl.serviceClasses></mp-object-icon>
                      </div>
                      <h2 class="card-pf-title text-center">
                        {{ service | objectName:$ctrl.serviceClasses }}
                      </h3>
                    </div>
                  </div>
                </div>
              </div>
            </div>`,
  bindings: {
    integrations: '<',
    serviceClasses: '<',
    serviceSelected: '&'
  }
});
