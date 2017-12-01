'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-integrations-list
 * @description
 * # mp-integrations-list
 */
angular.module('mobileControlPanelApp').component('mpIntegrationsList', {
  template: `<div ng-if="$ctrl.integrations.length === 0">
              <div class="empty-state-message text-center">
                <h2>ctrl Service doesn't have any Mobile Integrations available.</h2>
              </div>
            </div>

            <div ng-if="$ctrl.integrations.length > 0 && $ctrl.apps.length === 0">
              <div class="empty-state-message text-center">
                <h2>Get started with Mobile Integrations.</h2>
                <p class="gutter-top">
                Create a Mobile App to integrate with ctrl Service.
                </p>
                <p>
                  <a ng-href="project/{{ $ctrl.projectName }}/create-mobileapp" class="btn btn-primary btn-lg">Create Mobile App</a>
                </p>
              </div>
            </div>

            <div ng-if="$ctrl.integrations.length > 0 && $ctrl.apps.length > 0">
              <mp-service-integration ng-repeat="integration in $ctrl.integrations" integration=integration integration-toggled=$ctrl.integrationToggled></mp-service-integration>
            </div>`,
  bindings: {
    integrations: '<',
    integrationToggled: '<',
    apps: '<',
    projectName: '<'
  }
});
