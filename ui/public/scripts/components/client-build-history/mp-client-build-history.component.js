'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-client-build-history
 * @description
 * # mp-client-build-history
 */
angular.module('mobileControlPanelApp').component('mpClientBuildHistory', {
  template: `<mp-loading loading="$ctrl.loading">
              <mp-build-history
                build-config=$ctrl.buildConfig
                builds=$ctrl.builds>
              </mp-build-history>
            </mp-loading>`,
  controller: [
    '$scope',
    '$routeParams',
    'ClientBuildHistoryService',
    function($scope, $routeParams, ClientBuildHistoryService) {
      const watches = [];
      this.loading = true;

      ClientBuildHistoryService.getData($routeParams.project).then(data => {
        const [projectContext = {}, buildConfigs = {}, builds = {}] = data;

        this.projectContext = projectContext;

        const buildData = buildConfigs['_data'];
        this.buildConfig = Object.keys(buildData)
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

        this.builds = builds['_data'];

        watches.push(
          ClientBuildHistoryService.watch(
            'builds',
            this.projectContext,
            builds => {
              this.builds = Object.assign({}, builds['_data']);
            }
          )
        );

        this.loading = false;
      });

      $scope.$on('$destroy', function() {
        ClientBuildHistoryService.unwatchAll(watches);
      });
    }
  ]
});
