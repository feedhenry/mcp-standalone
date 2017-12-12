'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-mobile-client-tabs
 * @description
 * # mp-mobile-client-tabs
 */
angular.module('mobileControlPanelApp').component('mpMobileClientTabs', {
  template: `<div class="mp-object-tab-title">
              <h1>
                <span class="fa icon {{$ctrl.app.metadata.icon}} card-pf-icon-circle">&nbsp;</span>{{$ctrl.app.name}}
                <small class="meta" ng-if="$ctrl.app">created <span am-time-ago="$ctrl.app.metadata.created"></span></small>
              </h1>
            </div>
            <uib-tabset justified=true class="mp-tabsets" persist-tab-state>
              <uib-tab active="selectedTab.sdk" select="$ctrl.setEditorState($ctrl.states.VIEW)">
                <uib-tab-heading>SDK</uib-tab-heading>
                <mp-client-setup></mp-client-setup>
              </uib-tab>

              <uib-tab active="selectedTab.buildConfig">
                <uib-tab-heading>Build Config</uib-tab-heading>
                <mp-mobile-client-actions ng-if="$ctrl.buildConfig && $ctrl.hasMobileCiCd" action-selected=$ctrl.actionSelected></mp-mobile-client-actions>
                <mp-enable-mobile-build ng-if="!$ctrl.hasMobileCiCd"></mp-enable-mobile-build>
                <mp-client-build-config ng-if="$ctrl.hasMobileCiCd"</mp-client-build-config>
              </uib-tab>

              <uib-tab active="selectedTab.buildHistory" ng-if="$ctrl.buildConfig && $ctrl.hasMobileCiCd" select="$ctrl.setEditorState($ctrl.buildConfig ? $ctrl.states.VIEW : $ctrl.states.CREATE)">
                <uib-tab-heading>Build History</uib-tab-heading>
                <mp-mobile-client-actions ng-if="$ctrl.buildConfig && $ctrl.hasMobileCiCd" action-selected=$ctrl.actionSelected></mp-mobile-client-actions>
                <mp-client-build-history></mp-client-build-history>
              </uib-tab>
            </uib-tabset>`,
  bindings: {
    project: '<',
    app: '<',
    buildConfigs: '<',
    builds: '<',
    secrets: '<',
    startBuild: '&'
  },
  controller: [
    '$location',
    'ClientBuildEditorService',
    function($location, ClientBuildEditorService) {
      const MOBILE_CI_CD_NAME = 'aerogear-digger';
      this.states = ClientBuildEditorService.states;
      this.ClientBuildEditorService = ClientBuildEditorService;

      this.$onChanges = function(changes) {
        const buildConfigs =
          changes.buildConfigs && changes.buildConfigs.currentValue;
        if (buildConfigs) {
          const buildData = buildConfigs['_data'];
          this.buildConfig = Object.keys(buildData)
            .map(key => {
              return buildData[key];
            })
            .filter(buildConfig => {
              return (
                buildConfig.metadata.labels['mobile-appid'] === this.app.id
              );
            })
            .pop();
        }

        const secrets = changes.secrets && changes.secrets.currentValue;
        if (secrets) {
          this.hasMobileCiCd = Object.keys(secrets['_data'])
            .map(key => secrets['_data'][key])
            .some(secret => {
              return (
                secret.metadata.name === MOBILE_CI_CD_NAME &&
                secret.metadata.namespace === this.project.metadata.name
              );
            });
        }
      };

      this.setEditorState = function(state) {
        ClientBuildEditorService.state = state;
      };

      this.actionSelected = action => {
        if (action === 'build') {
          return this.startBuild()(this.buildConfig);
        }

        if (action === ClientBuildEditorService.states.EDIT) {
          ClientBuildEditorService.state = ClientBuildEditorService.states.EDIT;
          $location.url(
            `project/${this.project.metadata.name}/browse/mobileapps/${
              this.app.id
            }?tab=buildConfig`
          );
        }
      };
    }
  ]
});
