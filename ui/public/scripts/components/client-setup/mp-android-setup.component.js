'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-android-setup
 * @description
 * # mp-android-setup
 */
angular.module('mobileControlPanelApp').component('mpAndroidSetup', {
  template: `<div ng-if="$ctrl.app.clientType === 'android'">
              <h4>Installation</h4>
              <ul class="nav nav-pills nav-justified">
                <li role="presentation" ng-class="{'active' : $ctrl.installer === 'maven'}"><a ng-click="$ctrl.setInstaller('maven')" href="#">Maven</a></li>
                <li role="presentation" ng-class="{'active' : $ctrl.installer === 'gradle'}"><a ng-click="$ctrl.setInstaller('gradle')" href="#">Gradle</a></li>
              </ul>

              <p ng-if="$ctrl.installer === 'gradle'">Add the following to your <code>build.gradle</code> file</p>
              <p ng-if="$ctrl.installer === 'maven'">Add the following to your <code>pom.xml</code> file</p>

<mp-prettify ng-if="$ctrl.installer === 'maven'" type="'xml'" code-class="'prettyprint'">
&lt;dependency&gt;
  &lt;groupId&gt;org.feedhenry.mobile&lt;groupId&gt;
  &lt;artifactId&gt;mobile-core&lt;artifactId&gt;
  &lt;version&gt;0.0.1&lt;version&gt;
&lt;dependency&gt;
</mp-prettify>

<mp-prettify ng-if="$ctrl.installer === 'gradle'" type="'groovy'" code-class="'prettyprint'">
repositories {
  mavenCentral()
}
dependencies {
  compile group: 'org.feedhenry.mobile', name: 'mobile-core', version: '0.0.1'
}
</mp-prettify>

              <h4>Configuration</h4>
              <p ng-if="$ctrl.installer === 'gradle'">Add the following to a file named <code>mobile.properties</code> under your resources directory</p>
              <p ng-if="$ctrl.installer === 'maven'">Add the following to a file named <code>mobile.properties</code> under your resources directory</p>

<mp-prettify code-class="'prettyprint'">
org.feedhenry.mobile.host = "{{$ctrl.route}}"
org.feedhenry.mobile.appID = "{{$ctrl.app.id}}"
org.feedhenry.mobile.apiKey = "{{$ctrl.app.apiKey}}"
</mp-prettify>

              <h4>SDK Initialisation</h4>

<mp-prettify code-class="'prettyprint'" type="'java'">
import org.feedhenry.mobile
CoreSDK core = new CoreSDK()
try{
  CoreConfig cfg = cfgcore.configure(props)
}catch(InitialisationException e){
  //handle exception
}
</mp-prettify>

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
                    <tr>
                      <td>Hello World</td>
                      <td>Android hello world starter app. Shows you how to plugin the core sdk.</td>
                      <td><a href="https://github.com/feedhenry-templates/sync-android-app">https://github.com/feedhenry-templates/android-app</a></td>
                    </tr>
                    <tr>
                      <td>Sync Android quick start</td>
                      <td>A starting point for building out an application that syncs data automatically to the cloud</td>
                      <td><a href="https://github.com/feedhenry-templates/sync-android-app">https://github.com/feedhenry-templates/sync-android-app</a></td>
                    </tr>

                  </tbody>
                </table>
              </div>
            </div>`,
  bindings: {
    app: '<'
  },
  controller: [
    function() {
      this.installer = 'maven';
      this.route = window.MCP_URL;

      this.setInstaller = function(type) {
        this.installer = type;
      };
    }
  ]
});
