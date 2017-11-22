'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-ios-integration
 * @description
 * # mp-ios-integration
 */
angular.module('mobileControlPanelApp').component('mpIosIntegration', {
  template: `<h3>Installation</h3>

            <h3>Getting started</h3>

            <h3>Docs</h3>
            <ul>
              <!-- TODO: should this come from the service class? -->
              <li ng-if="$ctrl.service.name === 'fh-sync-server' "><a href="">Sync Server Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'fh-sync-server' "><a href="">Sync Client Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'keycloak' "><a href="">Keycloak Server Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'keycloak' "><a href="">Keycloak Client Documentation</a></li>
            </ul>

            <h3>Example Apps</h3>`,
  bindings: {
    service: '<'
  }
});
