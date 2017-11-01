'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-service-integration
 * @description
 * # mp-service-integration
 */
angular.module('mobileControlPanelApp').component('mpServiceIntegration', {
  template: `<div class="service-integration">
                <div class="integration-header row">
                  <div class="title">
                    <h3>{{$ctrl.integration.displayName}}</h3>
                  </div>
                  <div class="actions">
                    <mp-switch switched=switched disabled=!$ctrl.integration.target.writable checked=$ctrl.integration.enabled></mp-switch>
                    <span ng-if="!$ctrl.integration.target.writable">{{$ctrl.integration.target.displayName}} is not writeable; this means that it cannot be automatically configured with an integration.</span>
                  </div>
                </div>
                <div class="integration-info">
                  <div ng-include="'extensions/mcp/templates/' + $ctrl.integration.service + '-integration.template.html'"></div>
                </div>
              </div>`,
  bindings: {
    integration: '<',
    integrationToggled: '&'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.switched = function(event, value) {
        $scope.$ctrl.integrationToggled()($scope.$ctrl.integration, value);
      };
    }
  ]
});
