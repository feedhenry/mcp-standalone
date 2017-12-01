'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-build-history
 * @description
 * # mp-build-history
 */
angular.module('mobileControlPanelApp').component('mpBuildHistory', {
  template: `<div class="row">
                <div class="col-lg-12 build-trend">
                  <build-trends-chart $ctrl.builds=builds></build-trends-chart>
                </div>
              </div>
              <div class="row">
                <div class="col-lg-12">
                  <div class="build-pipelines" ng-repeat="build in $ctrl.orderedBuilds">
                    <build-pipeline build=build></build-pipeline>
                    <mp-app-download build=build></mp-app-download>
                  </div>
                </div>
              </div>`,
  bindings: {
    buildConfig: '<',
    builds: '<'
  },
  controller: [
    '$filter',
    'BuildsService',
    function($filter, BuildsService) {
      var buildConfigForBuild = $filter('buildConfigForBuild');
      this.$onChanges = function(changes) {
        const builds = changes.builds && changes.builds.currentValue;
        if (!builds) {
          return;
        }

        this.builds = _.filter(builds, build => {
          var buildConfigName = buildConfigForBuild(build) || '';
          return (
            this.buildConfig &&
            this.buildConfig.metadata.name === buildConfigName
          );
        });
        this.orderedBuilds = BuildsService.sortBuilds(this.builds, true);
      };
    }
  ]
});
