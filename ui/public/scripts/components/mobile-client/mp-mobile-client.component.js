'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-mobile-client
 * @description
 * # mp-mobile-client
 */
angular.module('mobileControlPanelApp').component('mpMobileClient', {
  template: `<div class="mp-app">
              <div class="container-fluid">
                <div class="row">
                  <div class="col-md-12">
                    <mp-mobile-client-tabs
                      project=$ctrl.project
                      app=$ctrl.app
                      build-configs=$ctrl.buildConfigs
                      builds=$ctrl.builds
                      secrets=$ctrl.secrets
                      start-build=$ctrl.startBuild>
                    </mp-mobile-client-tabs>
                  </div>
                </div>
              </div>
            </div>`,
  controller: [
    '$scope',
    '$location',
    '$routeParams',
    'MobileClientService',
    function($scope, $location, $routeParams, MobileClientService) {
      const watches = [];

      this.$onInit = function() {
        MobileClientService.getData(
          $routeParams.project,
          $routeParams.mobileapp
        ).then(data => {
          const [
            project = {},
            projectContext = {},
            app = {},
            buildConfigs = {},
            builds = {},
            secrets = {}
          ] = data;

          this.project = project;
          this.projectContext = projectContext;
          this.app = app;
          this.buildConfigs = buildConfigs;
          this.builds = builds;
          this.secrets = secrets;

          const watch = MobileClientService.watch(
            'buildconfigs',
            this.projectContext,
            buildConfigs => {
              this.buildConfigs = Object.assign({}, buildConfigs);
            }
          );
          watches.push(watch);
        });
      };

      this.startBuild = function(buildConfig) {
        MobileClientService.startBuild(buildConfig).then(() => {
          $location.url(
            `project/${$routeParams.project}/browse/mobileapps/${
              $routeParams.mobileapp
            }?tab=buildHistory`
          );
        });
      };

      $scope.$on('$destroy', function() {
        MobileClientService.unwatchAll(watches);
      });
    }
  ]
});
