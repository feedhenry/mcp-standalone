'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-client-build-config
 * @description
 * # mp-client-build-config
 */
angular.module('mobileControlPanelApp').component('mpClientBuildConfig', {
  template: `<mp-loading loading="$ctrl.loading">
                <mp-client-build-editor
                  build-config=$ctrl.buildConfig
                  config-created=$ctrl.createAppBuildConfig
                  config-updated=$ctrl.updateAppBuildConfig>
                </mp-client-build-editor>
              </mp-loading>`,
  controller: [
    '$routeParams',
    'ClientBuildConfigService',
    function($routeParams, ClientBuildConfigService) {
      const ctrl = this;
      ctrl.loading = true;

      ctrl.$onInit = function() {
        ClientBuildConfigService.getData($routeParams.project).then(data => {
          const [projectContext = {}, buildConfigs = []] = data;
          ctrl.projectContext = projectContext;
          ctrl.buildConfig = ctrl.getClientBuildConfig(buildConfigs['_data']);
          ctrl.loading = false;
        });
      };

      ctrl.createAppBuildConfig = function(appConfig) {
        appConfig.appID = $routeParams.mobileapp;
        ClientBuildConfigService.createBuildConfig(
          appConfig,
          ctrl.projectContext
        ).then(res => {
          ctrl.buildConfig = res;
        });
      };

      ctrl.updateAppBuildConfig = function(appConfig) {
        ClientBuildConfigService.updateBuildConfig(
          appConfig,
          ctrl.projectContext
        ).then(res => {
          ctrl.buildConfig = res;
        });
      };

      ctrl.getClientBuildConfig = function(buildData) {
        return Object.keys(buildData)
          .map(key => {
            return buildData[key];
          })
          .filter(buildConfig => {
            return (
              buildConfig.metadata.labels['mobile-appid'] ===
              $routeParams.mobileapp
            );
          })
          .pop();
      };
    }
  ]
});
