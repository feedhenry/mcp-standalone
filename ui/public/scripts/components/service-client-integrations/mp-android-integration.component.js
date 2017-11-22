'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-android-integration
 * @description
 * # mp-android-integration
 */
angular.module('mobileControlPanelApp').component('mpAndroidIntegration', {
  template: `<h3>Installation</h3>
            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && !$ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" type="'xml'" code-class="'indented-code prettyprint'">
              &lt;dependency&gt;
              &lt;groupId&gt;org.feedhenry.mobile&lt;/groupId&gt;
              &lt;artifactId&gt;mobile-core&lt;/artifactId&gt;
              &lt;version&gt;0.0.1&lt;/version&gt;
              &lt;/dependency&gt;
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && $ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" type="'xml'" code-class="'indented-code prettyprint'">
              &lt;dependency&gt;
              &lt;groupId&gt;org.feedhenry.mobile&lt;/groupId&gt;
              &lt;artifactId&gt;mobile-core&lt;/artifactId&gt;
              &lt;version&gt;0.0.1&lt;/version&gt;
              &lt;/dependency&gt;
              &lt;dependency&gt;
              &lt;groupId&gt;org.keycloak&lt;/groupId&gt;
              &lt;artifactId&gt;keycloak&lt;/artifactId&gt;
              &lt;version&gt;final-1.2.3&lt;/version&gt;
              &lt;/ dependency&gt;
            </mp-prettify>

            <h3>Getting started</h3>
            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && !$ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" code-class="'indented-code prettyprint'">
              android core sdk and sync code sample
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && $ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" code-class="'indented-code prettyprint'">
              android core sdk, sync, keycloak code sample
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && !$ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" code-class="'indented-code prettyprint'">
              sync client common use case functionality
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && $ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" code-class="'indented-code prettyprint'">
              sync client and keycloak common use case functionality
            </mp-prettify>

            <h3>Docs</h3>
            <ul>
              <!-- TODO: should this come from the service class? -->
              <li ng-if="$ctrl.service.name === 'fh-sync-server' "><a href="">Sync Server Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'fh-sync-server' "><a href="">Sync Client Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'keycloak' "><a href="">Keycloak Server Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'keycloak' "><a href="">Keycloak Client Documentation</a></li>
            </ul>

            <h3>Example Apps</h3>
            <div>
              <table class="table">
                <thead>
                  <tr>
                    <th>Template Name</th>
                    <th>Description</th>
                    <th>Source</th>
                  </tr>
                </thead>
                <tbody>
                  <tr ng-if="$ctrl.serviceClientService.enabled('keycloak', $ctrl.service)">
                    <td>Sync Integrated with Keycloak</td>
                    <td>Android app that integrates the core sdk and the keycloak and sync clients</td>
                    <td><a href="https://github.com/feedhenry-templates/sync-keycloak-app">https://github.com/feedhenry-templates/sync-keycloak-app</a></td>
                  </tr>
                  <tr >
                    <td>Sync Starter</td>
                    <td>Android starter app that integrates the core sdk and and sync clients</td>
                    <td><a href="https://github.com/feedhenry-templates/sync-app">https://github.com/feedhenry-templates/sync-app</a></td>
                  </tr>
                </tbody>
              </table>
            </div>`,
  bindings: {
    service: '<'
  },
  controller: [
    'ServiceClientService',
    function(ServiceClientService) {
      this.serviceClientService = ServiceClientService;
    }
  ]
});
