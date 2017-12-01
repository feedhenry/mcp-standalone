'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-client-integrations

 * @description
 * # mp-client-integrations
 */
angular.module('mobileControlPanelApp').component('mpClientIntegrations', {
  template: `<div ng-if="$ctrl.apps.length > 0">
              <ul ng-if="$ctrl.apps.length > 1" class="nav nav-pills nav-justified">
                <li ng-repeat="client in $ctrl.clients" role="presentation" ng-class="{'active' : client === $ctrl.clientType}"><a ng-click="$ctrl.installationOpt(client)" href="#">{{client}}</a></li>
              </ul>
              <p>
                Below you will find example code showing how to use the sdk client for {{$ctrl.service.name}}. You will also find sample code for use with a {{$ctrl.clientType}} app for any integrations that this service can take advantage of.
                You will also find some sample templates and quick starts to help you get started using your mobile service.
              </p>
              <mp-cordova-integration ng-if="$ctrl.clientType === 'cordova'" service=$ctrl.service></mp-cordova-integration>
              <mp-android-integration ng-if="$ctrl.clientType === 'android'" service=$ctrl.service></mp-android-integration>
              <mp-ios-integration ng-if="$ctrl.clientType === 'iOS'" service=$ctrl.service></mp-ios-integration>
            </div>`,
  bindings: {
    apps: '<',
    service: '<'
  },
  controller: [
    function() {
      this.$onChanges = function(changes) {
        const apps = changes.apps && changes.apps.currentValue;
        if (!apps) {
          return;
        }

        this.mobileapps = apps.reduce((acc, current) => {
          acc[current.clientType] = 'true';
          return acc;
        }, {});
        this.clients = Object.keys(this.mobileapps);
        this.clientType = this.clients[0] || 'cordoba';
      };

      this.installationOpt = function(type) {
        this.clientType = type;
      };
    }
  ]
});
