'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-sync-timings
 * @description
 * # mp-sync-timings
 */
angular.module('mobileControlPanelApp').component('mpSyncTimings', {
  template: `<div class="card-pf mp-sync-timings">
              <div class="card-pf-heading row">
                <h2 class="col-sm-6 col-xs-6">Timings (milliseconds)</h2>
                <span class="col-sm-6 col-xs-6">Last 30 Readings</span>
              </div>
              <div class="row">
                <mp-timings-chart chart-data=chartData></mp-timings-chart>
              </div>
            </div>`,
  bindings: {
    chartData: '<'
  },
  controller: [
    '$scope',
    function($scope) {
      const metricsToDisplay = ['mongodb_operation_time_', 'api_process_time_'];
      $scope.chartData = [];

      if (Array.isArray($scope.$ctrl.chartData[0])) {
        $scope.chartData.push($scope.$ctrl.chartData[0][0]);
      }

      $scope.$ctrl.chartData.forEach((chart, index) => {
        if (metricsToDisplay.indexOf(chart[1][0]) === -1) {
          return;
        }

        $scope.chartData.push(chart[1]);
      });
    }
  ]
});
