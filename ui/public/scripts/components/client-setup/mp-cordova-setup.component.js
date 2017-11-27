'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-cordova-setup
 * @description
 * # mp-cordova-setup
 */
angular.module('mobileControlPanelApp').component('mpCordovaSetup', {
  template: `<div ng-if="$ctrl.app.clientType === 'cordova'">
              <h4>Installation</h4>
              <p>Add the following to your <code>package.json</code> file</p>

<mp-prettify type="'bash'" code-class="'prettyprint'">
  npm install --save mobile-core
</mp-prettify>

              <h4>Configuration</h4>
              <p>Add the following to a file named <code>mobile.json</code> at the root of your project</p>
<mp-prettify type="'json'" code-class="'prettyprint'">
{
  "host":"{{route}}",
  "appID":"{{app.id}}",
  "apiKey":"{{app.apiKey}}"
}
</mp-prettify>

              <h4>SDK Initialisation</h4>
<mp-prettify type="'js'" code-class="'prettyprint'">
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
  });$ctrl.
});
</mp-prettify>

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
                    <tr>
                      <td>Hello World</td>
                      <td>Cordova hello world starter app. Shows you how to plugin the core sdk.</td>
                      <td><a href="https://github.com/feedhenry-templates/sync-cordova-app">https://github.com/feedhenry-templates/cordova-app</a>å</td>
                    </tr>
                    <tr>
                      <td>Sync Cordova quick start</td>
                      <td>A starting point for building out an application that syncs data automatically to the cloud</td>
                      <td><a href="https://github.com/feedhenry-templates/sync-cordova-app">https://github.com/feedhenry-templates/sync-cordova-app</a></td>
                    </tr>

                  </tbody>
                </table>
              </div>
            </div>`,
  bindings: {
    app: '<'
  }
});
