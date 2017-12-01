'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-cordova-integration
 * @description
 * # mp-cordova-integration
 */
angular.module('mobileControlPanelApp').component('mpCordovaIntegration', {
  template: `<h3>Installation</h3>
            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && !$ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" type="'bash'" code-class="'prettyprint indented-code'">
              npm install --save feedhenry-mobile-core-js
              npm install --save fh-sync-js@1.0.4
            </mp-prettify>
            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && $ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" type="'bash'" code-class="'prettyprint indented-code'">
              npm install --save fh-mobile-core
              npm install --save fh-sync-js@1.0.4
              npm install --save keycloak-js@1.0.4
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'keycloak'" type="'bash'" code-class="'prettyprint indented-code'">
                npm install --save fh-mobile-core
                npm install --save keycloak-js@1.0.4  
            </mp-prettify>

            <h3>Getting started</h3>
            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && !$ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" type="'js'" code-class="'prettyprint indented-code'">
              const mobileCore = require('fh-mobile-core');
              const sync = require('fh-sync-js');
              const mcpConfig = require('../mcpConfig.json');

              mobileCore.configure(mcpConfig).then((config) => {
                const syncConfig = config.getConfigFor('fh-sync-server');
                sync.init({
                  cloudUrl: syncConfig.uri,
                  storage_strategy: 'dom'
                });

                sync.manage('myDataset', null, {}, {}, () => {
                  // Initialise the rest of your app.
                });
              });
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'fh-sync-server' && $ctrl.serviceClientService.enabled('keycloak', $ctrl.service)" type="'js'" code-class="'prettyprint'">
              const mobileCore = require('fh-mobile-core');
              const sync = require('fh-sync-js');
              const keycloak = require('keycloak-js');
              const request = require('request-promise');
              const mcpConfig = require('../mcpConfig.json');

              function buildSyncCloudHandler(cloudUrl, options) {
                return function (params, success, failure) {
                  var url = cloudUrl + params.dataset_id;
                  var headers = (options.headers || {});
                  request({
                    method: 'POST',
                    uri: url,
                    headers: headers,
                    body: params.req,
                    json: true
                  })
                  .then((res) => {
                    return success(res);
                  })
                  .catch((err) => {
                    return failure(err);
                  });
                }
              }

              mobileCore.configure(mcpConfig).then((config) => {
                const syncConfig = config.getConfigFor('fh-sync-server');
                const keycloakConfig = config.getConfigFor('keycloak');

                sync.init({
                  cloudUrl: syncConfig.uri,
                  storage_strategy: 'dom'
                });

                sync.manage('myDataset', null, {}, {}, () => {
                  keycloak(keycloakConfig).init({ onLoad: 'login-required', flow: 'implicit' })
                  .success(() => {
                    const syncCloudUrl = syncConfig.uri + '/sync/';
                    const syncCloudHandler = buildSyncCloudHandler(syncCloudUrl, {
                      headers: {
                        'Authorization': 'Bearer ' + Keycloak.token
                      }
                    });
                    sync.setCloudHandler(syncCloudHandler);

                    // Initialise the rest of your app.
                  })
                });
              });
            </mp-prettify>

            <mp-prettify ng-if="$ctrl.service.name === 'keycloak'" type="'js'" code-class="'prettyprint'">
              const mobileCore = require('fh-mobile-core');
              const keycloak = require('keycloak-js');
              const mcpConfig = require('../mcpConfig.json');

              mobileCore.configure(mcpConfig).then((config) => {
                const keycloakConfig = config.getConfigFor('keycloak');

                keycloak(keycloakConfig).init({ onLoad: 'login-required', flow: 'implicit' })
                .success(() => {
                  // Handle Keycloak success and initialise the rest of your app.
                })
              })
            </mp-prettify>

            <h3>Docs</h3>
            <ul>
              <!-- TODO: should this come from the service class? -->
              <li ng-if="$ctrl.service.name === 'fh-sync-server'"><a href="">Sync Server Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'fh-sync-server'"><a href="">Sync Client Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'keycloak'"><a href="">Keycloak Server Documentation</a></li>
              <li ng-if="$ctrl.service.name === 'keycloak'"><a href="">Keycloak Client Documentation</a></li>
            </ul>

            <h3>Example Apps</h3>
            <!-- TODO: fix links to real docs & template, only if available. Otherwise leave empty -->
            <div>
              <!-- TODO: remove duplicate table header -->
              <table class="table">
                <thead>
                  <tr>
                    <th>Template Name</th>
                    <th>Description</th>
                    <th>Source</th>
                  </tr>
                </thead>
                <tbody>
                  <tr ng-if="$ctrl.service.name === 'fh-sync-server'">
                    <td>Sync Integrated with Keycloak</td>
                    <td>Cordova app that integrates the core sdk and the keycloak and sync clients</td>
                    <td><a href="https://github.com/feedhenry-templates/sync-keycloak-app">https://github.com/feedhenry-templates/sync-keycloak-app</a></td>
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
